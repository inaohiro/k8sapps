package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"k8soperation/pod/internal/models"
	"k8soperation/pod/internal/service"

	"github.com/go-chi/chi/v5"
)

func Controller() http.Handler {
	r := chi.NewRouter()
	r.Get("/", podIndex)
	r.Get("/{podID}", podDetail)
	r.Post("/", podCreate)
	r.Delete("/{podID}", podDelete)
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
	podID := chi.URLParam(r, "podID")
	pod, err := service.GetPod(r.Context(), namespace, podID)
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
	podID := chi.URLParam(r, "podID")
	err := service.DeletePod(r.Context(), namespace, podID)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Failed to delete pod: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
