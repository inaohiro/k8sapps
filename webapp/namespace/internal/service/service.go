package service

import (
	"context"
	"fmt"
	"k8soperation/core"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeleteNamespace(ctx context.Context, name string) error {
	clientset, err := core.GetKubeClient()
	if err != nil {
		return fmt.Errorf("Failed to create k8s client: %w", err)
	}

	clientset.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{})

	return nil
}
