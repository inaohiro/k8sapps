package controller

import (
	"encoding/json"
	"k8soperation/flavor/internal/service"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Controller() http.Handler {
	r := chi.NewRouter()
	r.Get("/", flavorIndex)
	return r
}

func flavorIndex(w http.ResponseWriter, r *http.Request) {
	flavors, err := service.ListFlavors(r.Context())
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flavors)
}
