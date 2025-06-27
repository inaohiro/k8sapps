package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v5"

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
				semconv.ServiceNameKey.String("webapp"),
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

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}
	_, err := strconv.Atoi(port)
	if err != nil {
		panic(fmt.Sprintf("failed to convert HTTP port number from HTTP_PORT environment variable into int: %v", err))
	}

	mux := http.NewServeMux()
	mux.Handle("POST /tokens", otelhttp.NewHandler(http.HandlerFunc(issueToken), "POST /tokens"))
	mux.Handle("GET /tokens", otelhttp.NewHandler(http.HandlerFunc(verifyToken), "GET /tokens"))

	http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}

type IssueTokenRequest struct {
	Namespace string `json:"namespace"`
}

type Token struct {
	Token string `json:"token"`
}

func issueToken(w http.ResponseWriter, r *http.Request) {
	var body IssueTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error("不正な値がリクエストボディにセットされています")
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"namespace": body.Namespace,
	})

	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		slog.Error("JWT トークンの署名に失敗しました")
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, Token{Token: tokenString})
}

const hmacSampleSecret = "hmacSampleSecret"

func verifyToken(w http.ResponseWriter, r *http.Request) {
	tokenString, err := getTokenFromAuthorizationHeader(r)
	if err != nil {
		slog.Error(err.Error())
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authorization header required"})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(hmacSampleSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		slog.Error("不正なトークンが渡されました")
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		slog.Error("不正なトークンが渡されました。claims の取得に失敗しました")
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid token"})
		return
	}

	namespace, ok := claims["namespace"].(string)
	if !ok {
		slog.Error("不正なトークンが渡されました。namespace の取得に失敗しました")
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid token"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"namespace": namespace})
}

func getTokenFromAuthorizationHeader(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return "", fmt.Errorf("不正な値がAuthorization headerにセットされています。expected: Bearer ${token}. got: %s", auth)
	}
	return auth[len("Bearer "):], nil
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error(err.Error())
	}
}
