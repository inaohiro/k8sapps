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
	r.Get("/", serviceIndex)
	r.Get("/{serviceID}", serviceDetail)
	r.Post("/", serviceCreate)
	r.Delete("/{serviceID}", serviceDelete)
	return r
}

func serviceIndex(w http.ResponseWriter, r *http.Request) {
	clientset, err := core.GetKubeClient()
	if err != nil {
		http.Error(w, "Failed to create k8s client: "+err.Error(), http.StatusInternalServerError)
		return
	}
	services, err := clientset.CoreV1().Services("").List(r.Context(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, "Failed to list services: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services.Items)
}

func serviceDetail(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "serviceID")
	clientset, err := core.GetKubeClient()
	if err != nil {
		http.Error(w, "Failed to create k8s client: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// serviceID is expected as namespace/name
	var namespace, name string
	_, err = fmt.Sscanf(serviceID, "%[^/]/%s", &namespace, &name)
	if err != nil {
		http.Error(w, "serviceID must be in namespace/name format", http.StatusBadRequest)
		return
	}
	service, err := clientset.CoreV1().Services(namespace).Get(r.Context(), name, metav1.GetOptions{})
	if err != nil {
		http.Error(w, "Failed to get service: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(service)
}

func serviceCreate(w http.ResponseWriter, r *http.Request) {
	clientset, err := core.GetKubeClient()
	if err != nil {
		http.Error(w, "Failed to create k8s client: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var serviceSpec map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&serviceSpec); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	serviceBytes, err := json.Marshal(serviceSpec)
	if err != nil {
		http.Error(w, "Failed to marshal service spec: "+err.Error(), http.StatusBadRequest)
		return
	}
	var serviceObj = &corev1.Service{}
	if err := json.Unmarshal(serviceBytes, serviceObj); err != nil {
		http.Error(w, "Failed to unmarshal service spec: "+err.Error(), http.StatusBadRequest)
		return
	}
	created, err := clientset.CoreV1().Services(serviceObj.Namespace).Create(r.Context(), serviceObj, metav1.CreateOptions{})
	if err != nil {
		http.Error(w, "Failed to create service: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func serviceDelete(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "serviceID")
	clientset, err := core.GetKubeClient()
	if err != nil {
		http.Error(w, "Failed to create k8s client: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var namespace, name string
	_, err = fmt.Sscanf(serviceID, "%[^/]/%s", &namespace, &name)
	if err != nil {
		http.Error(w, "serviceID must be in namespace/name format", http.StatusBadRequest)
		return
	}
	err = clientset.CoreV1().Services(namespace).Delete(r.Context(), name, metav1.DeleteOptions{})
	if err != nil {
		http.Error(w, "Failed to delete service: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
