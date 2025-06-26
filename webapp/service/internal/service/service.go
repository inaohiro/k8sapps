package service

import (
	"context"

	"k8soperation/core"
	"k8soperation/service/internal/models"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListServices(ctx context.Context, namespace string) ([]models.ServiceDTO, error) {
	clientset, err := core.GetKubeClient()
	if err != nil {
		return nil, err
	}
	services, err := clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return models.FromKubeServiceList(services.Items), nil
}

func GetService(ctx context.Context, namespace, name string) (*models.ServiceDTO, error) {
	clientset, err := core.GetKubeClient()
	if err != nil {
		return nil, err
	}
	svc, err := clientset.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	dto := models.FromKubeService(svc)
	return &dto, nil
}

func CreateService(ctx context.Context, namespace string, dto models.ServiceCreate) (*models.ServiceDTO, error) {
	clientset, err := core.GetKubeClient()
	if err != nil {
		return nil, err
	}
	svc := dto.ToKubeService()
	svc.Namespace = namespace
	created, err := clientset.CoreV1().Services(namespace).Create(ctx, svc, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	result := models.FromKubeService(created)
	return &result, nil
}

func DeleteService(ctx context.Context, namespace, name string) error {
	clientset, err := core.GetKubeClient()
	if err != nil {
		return err
	}
	return clientset.CoreV1().Services(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
