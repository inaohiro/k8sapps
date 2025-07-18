package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"k8soperation/core"
	"k8soperation/deployment/internal/models"
	"k8soperation/deployment/internal/service"

	"github.com/go-chi/chi/v5"
)

func Controller() http.Handler {
	r := chi.NewRouter()
	r.Get("/", deploymentIndex)
	r.Get("/{deploymentName}", deploymentDetail)
	r.Post("/", deploymentCreate)
	r.Delete("/{deploymentName}", deploymentDelete)
	return r
}

func deploymentIndex(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	deployments, err := service.ListDeployments(r.Context(), namespace)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	core.WriteJSON(w, http.StatusOK, deployments)
}

func deploymentDetail(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	deploymentName := chi.URLParam(r, "deploymentName")
	deployment, err := service.GetDeployment(r.Context(), namespace, deploymentName)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	core.WriteJSON(w, http.StatusOK, deployment)
}

func deploymentCreate(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	var req models.DeploymentCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	dto := models.FromRequest(req)
	created, err := service.CreateDeployment(r.Context(), namespace, dto)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	core.WriteJSON(w, http.StatusCreated, created)
}

func deploymentDelete(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	deploymentName := chi.URLParam(r, "deploymentName")
	if err := service.DeleteDeployment(r.Context(), namespace, deploymentName); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
