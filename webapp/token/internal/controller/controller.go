package controller

import (
	"encoding/json"
	"net/http"

	"k8soperation/token/internal/models"
	"k8soperation/token/internal/service"

	"github.com/go-chi/chi/v5"
)

func Controller() http.Handler {
	r := chi.NewRouter()
	r.Post("/", tokenIssue)
	return r
}

func tokenIssue(w http.ResponseWriter, r *http.Request) {
	var body models.IssueTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	token, err := service.IssueToken(models.FromRequest(body))
	if err != nil {
		http.Error(w, "Failed to issue token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Token{
		Token: token,
	})
}
