package service

import (
	"context"

	"k8soperation/core"

	"k8soperation/pod/internal/models"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListPods returns a list of pods in the given namespace
func ListPods(ctx context.Context, namespace string) ([]models.Pod, error) {
	clientset, err := core.GetKubeClient()
	if err != nil {
		return nil, err
	}
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return models.FromKubePodList(pods.Items), nil
}

// GetPod returns a pod by name in the given namespace
func GetPod(ctx context.Context, namespace, podID string) (*models.Pod, error) {
	clientset, err := core.GetKubeClient()
	if err != nil {
		return nil, err
	}
	pod, err := clientset.CoreV1().Pods(namespace).Get(ctx, podID, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	dto := models.FromKubePod(pod)
	return &dto, nil
}

// CreatePod creates a new pod from a PodCreate DTO
func CreatePod(ctx context.Context, namespace string, pod models.PodCreate) (*models.Pod, error) {
	clientset, err := core.GetKubeClient()
	if err != nil {
		return nil, err
	}
	podObj := pod.ToKubePod()
	podObj.Namespace = namespace
	created, err := clientset.CoreV1().Pods(namespace).Create(ctx, podObj, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	dto := models.FromKubePod(created)
	return &dto, nil
}

// DeletePod deletes a pod by name in the given namespace
func DeletePod(ctx context.Context, namespace, podID string) error {
	clientset, err := core.GetKubeClient()
	if err != nil {
		return err
	}
	return clientset.CoreV1().Pods(namespace).Delete(ctx, podID, metav1.DeleteOptions{})
}
