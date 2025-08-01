package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"

	_ "embed"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

//go:embed config.json
var configjson string

var (
	auth_url string
	app_url  string
)

type Route struct {
	Pattern     string `json:"pattern"`
	IssueToken  bool   `json:"issue_token"`
	VerifyToken bool   `json:"verify_token"`
	Proxy       bool   `json:"proxy"`
}

var (
	initResourcesOnce sync.Once
	resource          *sdkresource.Resource
)

func initResource() *sdkresource.Resource {
	initResourcesOnce.Do(func() {
		extraResources, _ := sdkresource.New(
			context.Background(),
			sdkresource.WithOS(),
			sdkresource.WithProcess(),
			sdkresource.WithContainer(),
			sdkresource.WithHost(),
			sdkresource.WithAttributes(
				semconv.ServiceNameKey.String("gateway"),
				semconv.ServiceVersionKey.String("1.0.0"),
			),
		)
		resource, _ = sdkresource.Merge(
			sdkresource.Default(),
			extraResources,
		)
	})
	return resource
}

func initTracerProvider() *sdktrace.TracerProvider {
	ctx := context.Background()

	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		log.Fatalf("OTLP Trace gRPC Creation: %v", err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(initResource()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func initMeterProvider() *sdkmetric.MeterProvider {
	ctx := context.Background()

	exporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		log.Fatalf("new otlp metric grpc exporter failed: %v", err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(initResource()),
	)
	otel.SetMeterProvider(mp)
	return mp
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	tp := initTracerProvider()
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("Tracer Provider Shutdown: %v", err)
		}
		log.Println("Shutdown tracer provider")
	}()

	mp := initMeterProvider()
	defer func() {
		if err := mp.Shutdown(ctx); err != nil {
			log.Fatalf("Error shutting down meter provider: %v", err)
		}
		log.Println("Shutdown meter provider")
	}()

	var routes []Route
	if err := json.Unmarshal([]byte(configjson), &routes); err != nil {
		log.Fatal(err)
	}

	var err error
	auth_url = os.Getenv("AUTH_URL")
	if auth_url == "" {
		log.Fatal("AUTH_URL required")
	}
	_, err = url.Parse(auth_url)
	if err != nil {
		log.Fatal(err)
	}
	app_url = os.Getenv("APP_URL")
	if app_url == "" {
		log.Fatal("APP_URL required")
	}
	_, err = url.Parse(app_url)
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}
	_, err = strconv.Atoi(port)
	if err != nil {
		log.Fatalf("failed to convert HTTP port number from HTTP_PORT environment variable into int: %v", err)
	}

	mux := http.NewServeMux()
	for _, route := range routes {
		handler := nopHandler()
		if route.Proxy {
			handler = proxy(handler)
		}
		if route.VerifyToken {
			handler = tokenVerify(handler)
		}
		if route.IssueToken {
			handler = issueToken(handler)
		}

		mux.Handle(route.Pattern, otelhttp.NewHandler(handler, route.Pattern))
	}

	http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}

var (
	client = http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
)

type VerifyTokenResponse struct {
	Namespace string `json:"namespace"`
	Error     string `json:"error"`
}

type contextKeyNamespace struct{}

func tokenVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// トークン検証
		endpoint, err := url.JoinPath(auth_url, "tokens")
		if err != nil {
			slog.Error("failed to url.JoinPath", slog.String("err", err.Error()))
			writeJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
			return
		}
		req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, strings.TrimSpace(endpoint), nil)
		if err != nil {
			slog.Error("failed to http.NewRequestWithContext", slog.String("err", err.Error()))
			writeJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
			return
		}
		// Authorization ヘッダが付いているはず
		// トークン検証 API につけて送信
		req.Header = r.Header
		resp, err := client.Do(req)
		if err != nil {
			slog.Error("failed to send request", slog.String("err", err.Error()))
			writeJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
			return
		}
		defer resp.Body.Close()
		_body, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("failed to read response body", slog.String("err", err.Error()))
			writeJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
			return
		}
		var verifyTokenResponse VerifyTokenResponse
		if err := json.Unmarshal(_body, &verifyTokenResponse); err != nil {
			slog.Error("failed to unmarshal response body", slog.String("err", err.Error()))
			writeJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
			return
		}
		if resp.StatusCode >= 400 {
			slog.Error("トークン検証に失敗しました", slog.String("err", verifyTokenResponse.Error))
			writeJSON(w, resp.StatusCode, map[string]string{"message": verifyTokenResponse.Error})
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKeyNamespace{}, verifyTokenResponse.Namespace)))
	})
}

func proxy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		// proxy するには token 検証が必要
		namespace, ok := r.Context().Value(contextKeyNamespace{}).(string)

		// トークン検証が成功したらアプリケーションにリクエスト送信
		// 元のリクエストパスに /api/{namespace} をつける
		var endpoint string
		if ok {
			var err error
			endpoint, err = url.JoinPath(app_url, "api", namespace, strings.Join(strings.Split(r.URL.Path, "/")[2:], "/"))
			if err != nil {
				slog.Error("failed to url.JoinPath", slog.String("err", err.Error()))
				writeJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
				return
			}
		} else {
			// トークン検証なしで proxy が実行されることもある
			// リクエストパスをそのまま使う
			var err error
			endpoint, err = url.JoinPath(app_url, r.URL.Path)
			if err != nil {
				slog.Error("failed to url.JoinPath", slog.String("err", err.Error()))
				writeJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
				return
			}
		}

		req, err := http.NewRequestWithContext(r.Context(), r.Method, endpoint, r.Body)
		if err != nil {
			slog.Error("failed to http.NewRequestWithContext", slog.String("err", err.Error()))
			writeJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
			return
		}
		req.Header = r.Header
		resp, err := client.Do(req)
		if err != nil {
			slog.Error("failed to send request", slog.String("err", err.Error()))
			writeJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
			return
		}

		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})
}

func issueToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		// トークン発行
		endpoint, err := url.JoinPath(auth_url, "tokens")
		if err != nil {
			slog.Error("failed to url.JoinPath", slog.String("err", err.Error()))
			writeJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
			return
		}
		req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, endpoint, r.Body)
		if err != nil {
			slog.Error("failed to http.NewRequestWithContext", slog.String("err", err.Error()))
			writeJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
			return
		}
		req.Header = r.Header
		resp, err := client.Do(req)
		if err != nil {
			slog.Error("failed to send request", slog.String("err", err.Error()))
			writeJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
			return
		}

		// そのままレスポンスを返す
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})
}

func nopHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error(err.Error())
	}
}
