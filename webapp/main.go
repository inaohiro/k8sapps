package main

import (
	"context"
	"fmt"
	"k8soperation/core"
	mymiddleware "k8soperation/core/middleware"
	"k8soperation/deployment"
	"k8soperation/flavor"
	"k8soperation/image"
	"k8soperation/namespace"
	"k8soperation/pod"
	"k8soperation/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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

	core.InitDB()

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

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(mymiddleware.CreateNamespace)

	patternRouteMap := map[string]http.Handler{
		"/api/{namespace}/pods":        pod.Routes,
		"/api/{namespace}/deployments": deployment.Routes,
		"/api/{namespace}/services":    service.Routes,
		"/api/images":                  image.Routes,
		"/api/flavors":                 flavor.Routes,
		"/api/namespace":               namespace.Routes,
	}
	for pattern, handler := range patternRouteMap {
		r.Mount(pattern, otelhttp.NewHandler(handler, pattern))
	}

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}
	_, err := strconv.Atoi(port)
	if err != nil {
		panic(fmt.Sprintf("failed to convert HTTP port number from HTTP_PORT environment variable into int: %v", err))
	}

	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
