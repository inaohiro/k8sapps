package controller

import (
	"encoding/json"
	"k8soperation/image/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Controller() http.Handler {
	r := chi.NewRouter()
	r.Get("/", imageIndex)
	return r
}

func imageIndex(w http.ResponseWriter, r *http.Request) {
	images, err := service.ListImages(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(images)
}
