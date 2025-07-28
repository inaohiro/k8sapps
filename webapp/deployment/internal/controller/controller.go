package controller

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"k8soperation/core"
	"k8soperation/core/middleware"
	"k8soperation/deployment/internal/models"
	"k8soperation/deployment/internal/service"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func Controller() http.Handler {
	r := chi.NewRouter()
	prefix := "/api/{namespace}/deployments"

	r.Method(http.MethodGet, "/", otelhttp.NewHandler(middleware.IntentionalError(http.HandlerFunc(deploymentIndex)), fmt.Sprintf("GET %s%s", prefix, "")))
	r.Method(http.MethodGet, "/{deployment_name}", otelhttp.NewHandler(middleware.IntentionalError(http.HandlerFunc(deploymentDetail)), fmt.Sprintf("GET %s%s", prefix, "/{deployment_name}")))
	r.Method(http.MethodPost, "/", otelhttp.NewHandler(middleware.IntentionalError(http.HandlerFunc(deploymentCreate)), fmt.Sprintf("POST %s%s", prefix, "")))
	r.Method(http.MethodDelete, "/{deployment_name}", otelhttp.NewHandler(middleware.IntentionalError(http.HandlerFunc(deploymentDelete)), fmt.Sprintf("DELETE %s%s", prefix, "/{deployment_name}")))
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
	name := chi.URLParam(r, "deployment_name")
	deployment, err := service.GetDeployment(r.Context(), namespace, name)
	if err != nil {
		status := core.GetErrorStatus(err)
		http.Error(w, err.Error(), status)
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
	name := chi.URLParam(r, "deployment_name")
	if err := service.DeleteDeployment(r.Context(), namespace, name); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
