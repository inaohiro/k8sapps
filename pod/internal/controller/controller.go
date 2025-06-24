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
	r.Get("/", podIndex)
	r.Get("/{podID}", podDetail)
	r.Post("/", podCreate)
	r.Delete("/{podID}", podDelete)
	return r
}

func podIndex(w http.ResponseWriter, r *http.Request) {
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

func podDetail(w http.ResponseWriter, r *http.Request) {
	podID := chi.URLParam(r, "podID")
	clientset, err := core.GetKubeClient()
	if err != nil {
		http.Error(w, "Failed to create k8s client: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// podID is expected as namespace/name
	var namespace, name string
	_, err = fmt.Sscanf(podID, "%[^/]/%s", &namespace, &name)
	if err != nil {
		http.Error(w, "podID must be in namespace/name format", http.StatusBadRequest)
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

func podCreate(w http.ResponseWriter, r *http.Request) {
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

func podDelete(w http.ResponseWriter, r *http.Request) {
	podID := chi.URLParam(r, "podID")
	clientset, err := core.GetKubeClient()
	if err != nil {
		http.Error(w, "Failed to create k8s client: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var namespace, name string
	_, err = fmt.Sscanf(podID, "%[^/]/%s", &namespace, &name)
	if err != nil {
		http.Error(w, "podID must be in namespace/name format", http.StatusBadRequest)
		return
	}
	err = clientset.CoreV1().Pods(namespace).Delete(r.Context(), name, metav1.DeleteOptions{})
	if err != nil {
		http.Error(w, "Failed to delete pod: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
