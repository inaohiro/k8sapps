package controller

import (
	"encoding/json"
	"net/http"

	"k8soperation/core"
	"k8soperation/service/internal/models"
	"k8soperation/service/internal/service"

	"github.com/go-chi/chi/v5"
)

func Controller() http.Handler {
	r := chi.NewRouter()
	r.Get("/", serviceIndex)
	r.Get("/{serviceID}", serviceDetail)
	r.Post("/", serviceCreate)
	r.Delete("/{serviceID}", serviceDelete)
	return r
}

func serviceIndex(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	services, err := service.ListServices(r.Context(), namespace)
	if err != nil {
		http.Error(w, "Failed to list services: "+err.Error(), http.StatusInternalServerError)
		return
	}
	core.WriteJSON(w, http.StatusOK, services)
}

func serviceDetail(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	serviceID := chi.URLParam(r, "serviceID")
	svc, err := service.GetService(r.Context(), namespace, serviceID)
	if err != nil {
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
		http.Error(w, "Failed to create service: "+err.Error(), http.StatusInternalServerError)
		return
	}
	core.WriteJSON(w, http.StatusCreated, created)
}

func serviceDelete(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	serviceID := chi.URLParam(r, "serviceID")
	err := service.DeleteService(r.Context(), namespace, serviceID)
	if err != nil {
		http.Error(w, "Failed to delete service: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
