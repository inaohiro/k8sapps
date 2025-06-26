package service

import (
	"context"
	"encoding/json"

	"k8soperation/core"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListPods returns a list of pods in the given namespace
func ListPods(ctx context.Context, namespace string) ([]corev1.Pod, error) {
	clientset, err := core.GetKubeClient()
	if err != nil {
		return nil, err
	}
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return pods.Items, nil
}

// GetPod returns a pod by name in the given namespace
func GetPod(ctx context.Context, namespace, podID string) (*corev1.Pod, error) {
	clientset, err := core.GetKubeClient()
	if err != nil {
		return nil, err
	}
	pod, err := clientset.CoreV1().Pods(namespace).Get(ctx, podID, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return pod, nil
}

// CreatePod creates a new pod from a generic map spec
func CreatePod(ctx context.Context, podSpec map[string]interface{}) (*corev1.Pod, error) {
	clientset, err := core.GetKubeClient()
	if err != nil {
		return nil, err
	}
	podBytes, err := json.Marshal(podSpec)
	if err != nil {
		return nil, err
	}
	var podObj = &corev1.Pod{}
	if err := json.Unmarshal(podBytes, podObj); err != nil {
		return nil, err
	}
	created, err := clientset.CoreV1().Pods(podObj.Namespace).Create(ctx, podObj, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return created, nil
}

// DeletePod deletes a pod by name in the given namespace
func DeletePod(ctx context.Context, namespace, podID string) error {
	clientset, err := core.GetKubeClient()
	if err != nil {
		return err
	}
	return clientset.CoreV1().Pods(namespace).Delete(ctx, podID, metav1.DeleteOptions{})
}
