package controller

import (
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
	r.Method(http.MethodDelete, "/{namespace_name}", otelhttp.NewHandler(
		middleware.IntentionalError(http.HandlerFunc(namespaceDelete)),
		fmt.Sprintf("DELETE %s%s", prefix, "/{namespace_name}"),
	))

	return r
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
