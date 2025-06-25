package otel

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceProvider := otel.GetTracerProvider()
		tracer := traceProvider.Tracer("isucon14", oteltrace.WithInstrumentationVersion("0.1.0"))

		opts := []oteltrace.SpanStartOption{
			oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", r)...),
			oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(r)...),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		}
		ctx, span := tracer.Start(r.Context(), "", opts...)
		defer span.End()
		myw := &myResponseWriter{status: 200, w: w}

		next.ServeHTTP(myw, r.WithContext(ctx))

		c := chi.RouteContext(ctx)
		span.SetName(c.RoutePattern())
		span.SetAttributes(semconv.HTTPServerAttributesFromHTTPRequest("isucon14", c.RoutePattern(), r)...)
		span.SetAttributes(semconv.HTTPAttributesFromHTTPStatusCode(myw.status)...)

		spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(myw.status)
		span.SetStatus(spanStatus, spanMessage)
	})
}

type myResponseWriter struct {
	status int
	w      http.ResponseWriter
}

func (w *myResponseWriter) Header() http.Header {
	return w.w.Header()
}
func (w *myResponseWriter) Write(b []byte) (int, error) {
	return w.w.Write(b)
}
func (w *myResponseWriter) WriteHeader(statusCode int) {
	w.w.WriteHeader(statusCode)
	w.status = statusCode
}
