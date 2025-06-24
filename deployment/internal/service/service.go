package service

import (
	"context"
	"fmt"

	"k8soperation/core"
	"k8soperation/deployment/internal/models"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListDeployments returns all deployments in all namespaces
func ListDeployments(ctx context.Context) ([]appsv1.Deployment, error) {
	clientset, err := core.GetKubeClient()
	if err != nil {
		return nil, fmt.Errorf("Failed to create k8s client: %w", err)
	}
	deployments, err := clientset.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Failed to list deployments: %w", err)
	}
	return deployments.Items, nil
}

// GetDeployment returns a deployment by namespace and name
func GetDeployment(ctx context.Context, namespace, name string) (*appsv1.Deployment, error) {
	clientset, err := core.GetKubeClient()
	if err != nil {
		return nil, fmt.Errorf("Failed to create k8s client: %w", err)
	}
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("Failed to get deployment: %w", err)
	}
	return deployment, nil
}

// CreateDeployment creates a deployment from a DeploymentCreate DTO
func CreateDeployment(ctx context.Context, namespace string, dto models.DeploymentCreate) (*appsv1.Deployment, error) {
	clientset, err := core.GetKubeClient()
	if err != nil {
		return nil, fmt.Errorf("Failed to create k8s client: %w", err)
	}
	deploymentObj := dto.ToKubeDeployment()
	created, err := clientset.AppsV1().Deployments(namespace).Create(ctx, deploymentObj, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("Failed to create deployment: %w", err)
	}
	return created, nil
}

// DeleteDeployment deletes a deployment by namespace and name
func DeleteDeployment(ctx context.Context, namespace, name string) error {
	clientset, err := core.GetKubeClient()
	if err != nil {
		return fmt.Errorf("Failed to create k8s client: %w", err)
	}
	if err := clientset.AppsV1().Deployments(namespace).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("Failed to delete deployment: %w", err)
	}
	return nil
}
