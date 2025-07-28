package controller

import (
	"encoding/json"
	"fmt"
	"k8soperation/core/middleware"
	"k8soperation/image/internal/service"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func Controller() http.Handler {
	r := chi.NewRouter()
	prefix := "/api/images"
	r.Method(http.MethodGet, "/", otelhttp.NewHandler(
		middleware.IntentionalError(http.HandlerFunc(imageIndex)),
		fmt.Sprintf("GET %s%s", prefix, ""),
	))
	return r
}

func imageIndex(w http.ResponseWriter, r *http.Request) {
	images, err := service.ListImages(r.Context())
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(images)
}
