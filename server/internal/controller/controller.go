package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"k8soperation/core"

	"github.com/go-chi/chi/v5"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Controller() http.Handler {
	r := chi.NewRouter()
	r.Get("/", serverIndex)
	r.Get("/{serverID}", serverDetail)
	r.Post("/", serverCreate)
	r.Delete("/{serverID}", serverDelete)
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

func serverCreate(w http.ResponseWriter, r *http.Request) {
	clientset, err := core.GetKubeClient()
	if err != nil {
		http.Error(w, "Failed to create k8s client: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var podSpec map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&podSpec); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Marshal and unmarshal to concrete type
	podBytes, err := json.Marshal(podSpec)
	if err != nil {
		http.Error(w, "Failed to marshal pod spec: "+err.Error(), http.StatusBadRequest)
		return
	}
	var podObj = &corev1.Pod{}
	if err := json.Unmarshal(podBytes, podObj); err != nil {
		http.Error(w, "Failed to unmarshal pod spec: "+err.Error(), http.StatusBadRequest)
		return
	}
	created, err := clientset.CoreV1().Pods(podObj.Namespace).Create(r.Context(), podObj, metav1.CreateOptions{})
	if err != nil {
		http.Error(w, "Failed to create pod: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func serverDelete(w http.ResponseWriter, r *http.Request) {
	serverID := chi.URLParam(r, "serverID")
	clientset, err := core.GetKubeClient()
	if err != nil {
		http.Error(w, "Failed to create k8s client: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var namespace, name string
	_, err = fmt.Sscanf(serverID, "%[^/]/%s", &namespace, &name)
	if err != nil {
		http.Error(w, "serverID must be in namespace/name format", http.StatusBadRequest)
		return
	}
	err = clientset.CoreV1().Pods(namespace).Delete(r.Context(), name, metav1.DeleteOptions{})
	if err != nil {
		http.Error(w, "Failed to delete pod: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
