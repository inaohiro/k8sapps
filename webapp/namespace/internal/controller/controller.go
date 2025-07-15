package controller

import (
	"k8soperation/core"
	"k8soperation/namespace/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Controller() http.Handler {
	r := chi.NewRouter()
	r.Delete("/{namespaceName}", namespaceDelete)

	return r
}

func namespaceDelete(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "namespaceName")
	err := service.DeleteNamespace(r.Context(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	core.WriteJSON(w, http.StatusNoContent, nil)
}
