package controller

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"k8soperation/core"
	"k8soperation/core/middleware"
	"k8soperation/service/internal/models"
	"k8soperation/service/internal/service"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func Controller() http.Handler {
	r := chi.NewRouter()
	prefix := "/api/{namespace}/services"

	r.Method(http.MethodGet, "/", otelhttp.NewHandler(middleware.IntentionalError(http.HandlerFunc(serviceIndex)), fmt.Sprintf("GET %s%s", prefix, "")))
	r.Method(http.MethodGet, "/{service_name}", otelhttp.NewHandler(middleware.IntentionalError(http.HandlerFunc(serviceDetail)), fmt.Sprintf("GET %s%s", prefix, "/{service_name}")))
	r.Method(http.MethodPost, "/", otelhttp.NewHandler(middleware.IntentionalError(http.HandlerFunc(serviceCreate)), fmt.Sprintf("POST %s%s", prefix, "")))
	r.Method(http.MethodDelete, "/{service_name}", otelhttp.NewHandler(middleware.IntentionalError(http.HandlerFunc(serviceDelete)), fmt.Sprintf("DELETE %s%s", prefix, "/{service_name}")))
	return r
}

func serviceIndex(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	services, err := service.ListServices(r.Context(), namespace)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Failed to list services: "+err.Error(), http.StatusInternalServerError)
		return
	}
	core.WriteJSON(w, http.StatusOK, services)
}

func serviceDetail(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	serviceID := chi.URLParam(r, "service_name")
	svc, err := service.GetService(r.Context(), namespace, serviceID)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Failed to get service: "+err.Error(), http.StatusInternalServerError)
		return
	}
	core.WriteJSON(w, http.StatusOK, svc)
}

func serviceCreate(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	var svcObj models.ServiceCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&svcObj); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	dto := models.FromRequest(svcObj)
	created, err := service.CreateService(r.Context(), namespace, dto)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Failed to create service: "+err.Error(), http.StatusInternalServerError)
		return
	}
	core.WriteJSON(w, http.StatusCreated, created)
}

func serviceDelete(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	serviceID := chi.URLParam(r, "service_name")
	err := service.DeleteService(r.Context(), namespace, serviceID)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Failed to delete service: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
