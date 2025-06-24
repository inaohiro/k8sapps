package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Controller() http.Handler {
	r := chi.NewRouter()
	r.Get("/", imageIndex)
	return r
}

func imageIndex(w http.ResponseWriter, r *http.Request) {

}
