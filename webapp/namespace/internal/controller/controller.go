package controller

import (
	"encoding/json"
	"fmt"
	"k8soperation/core/middleware"
	"k8soperation/namespace/internal/service"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func Controller() http.Handler {
	r := chi.NewRouter()
	prefix := "/api/namespace"

	r.Method(http.MethodGet, "/", otelhttp.NewHandler(
		middleware.IntentionalError(http.HandlerFunc(namespaceIndex)),
		fmt.Sprintf("GET %s", prefix),
	))
	r.Method(http.MethodDelete, "/_all", otelhttp.NewHandler(
		middleware.IntentionalError(http.HandlerFunc(namespaceDeleteAll)),
		fmt.Sprintf("DELETE %s%s", prefix, "/_all"),
	))
	r.Method(http.MethodDelete, "/{namespace_name}", otelhttp.NewHandler(
		middleware.IntentionalError(http.HandlerFunc(namespaceDelete)),
		fmt.Sprintf("DELETE %s%s", prefix, "/{namespace_name}"),
	))

	return r
}

func namespaceIndex(w http.ResponseWriter, r *http.Request) {
	namespaces, err := service.ListNamespaces(r.Context())
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(namespaces)
}

func namespaceDelete(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "namespace_name")
	err := service.DeleteNamespace(r.Context(), name)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func namespaceDeleteAll(w http.ResponseWriter, r *http.Request) {
	err := service.DeleteAllNamespaces(r.Context())
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
