package controller

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"k8soperation/core/middleware"
	"k8soperation/pod/internal/models"
	"k8soperation/pod/internal/service"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func Controller() http.Handler {
	r := chi.NewRouter()
	prefix := "/api/{namespace}/pods"

	r.Method(http.MethodGet, "/", otelhttp.NewHandler(middleware.IntentionalError(http.HandlerFunc(podIndex)), fmt.Sprintf("GET %s%s", prefix, "")))
	r.Method(http.MethodGet, "/{pod_name}", otelhttp.NewHandler(middleware.IntentionalError(http.HandlerFunc(podDetail)), fmt.Sprintf("GET %s%s", prefix, "/{pod_name}")))
	r.Method(http.MethodPost, "/", otelhttp.NewHandler(middleware.IntentionalError(http.HandlerFunc(podCreate)), fmt.Sprintf("POST %s%s", prefix, "")))
	r.Method(http.MethodDelete, "/{pod_name}", otelhttp.NewHandler(middleware.IntentionalError(http.HandlerFunc(podDelete)), fmt.Sprintf("DELETE %s%s", prefix, "/{pod_name}")))

	return r
}

func podIndex(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	pods, err := service.ListPods(r.Context(), namespace)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Failed to list pods: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pods)
}

func podDetail(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	name := chi.URLParam(r, "pod_name")
	pod, err := service.GetPod(r.Context(), namespace, name)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Failed to get pod: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pod)
}

func podCreate(w http.ResponseWriter, r *http.Request) {
	var req models.PodCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	podDTO := models.FromRequest(req)
	namespace := chi.URLParam(r, "namespace")
	created, err := service.CreatePod(r.Context(), namespace, podDTO)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Failed to create pod: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func podDelete(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	name := chi.URLParam(r, "pod_name")
	err := service.DeletePod(r.Context(), namespace, name)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Failed to delete pod: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
