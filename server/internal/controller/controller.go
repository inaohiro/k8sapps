package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"inpodk8soperation/core"

	"github.com/go-chi/chi/v5"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Controller() http.Handler {
	r := chi.NewRouter()
	r.Get("/", serverIndex)
	r.Get("/{serverID}", serverDetail)
	return r
}

func serverIndex(w http.ResponseWriter, r *http.Request) {
	clientset, err := core.GetKubeClient()
	if err != nil {
		http.Error(w, "Failed to create k8s client: "+err.Error(), http.StatusInternalServerError)
		return
	}
	pods, err := clientset.CoreV1().Pods("").List(r.Context(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, "Failed to list pods: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pods.Items)
}

func serverDetail(w http.ResponseWriter, r *http.Request) {
	serverID := chi.URLParam(r, "serverID")
	clientset, err := core.GetKubeClient()
	if err != nil {
		http.Error(w, "Failed to create k8s client: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// serverID is expected as namespace/name
	var namespace, name string
	_, err = fmt.Sscanf(serverID, "%[^/]/%s", &namespace, &name)
	if err != nil {
		http.Error(w, "serverID must be in namespace/name format", http.StatusBadRequest)
		return
	}
	pod, err := clientset.CoreV1().Pods(namespace).Get(r.Context(), name, metav1.GetOptions{})
	if err != nil {
		http.Error(w, "Failed to get pod: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pod)
}
