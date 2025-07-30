package service

import (
	"context"
	"fmt"
	"k8soperation/core"
	"strings"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListNamespaces(ctx context.Context) ([]string, error) {
	clientset, err := core.GetKubeClient()
	if err != nil {
		return nil, fmt.Errorf("Failedt to create k8s client: %w", err)
	}

	namespaceList, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(namespaceList.Items))
	for _, v := range namespaceList.Items {
		if strings.HasPrefix(v.Name, "k8sapps") {
			result = append(result, v.Name)
		}
	}

	return result, nil

}

func DeleteNamespace(ctx context.Context, name string) error {
	clientset, err := core.GetKubeClient()
	if err != nil {
		return fmt.Errorf("Failed to create k8s client: %w", err)
	}

	clientset.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{})

	return nil
}

func DeleteAllNamespaces(ctx context.Context) error {
	namespaces, err := ListNamespaces(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(namespaces))
	for _, name := range namespaces {
		go func() {
			defer wg.Done()
			DeleteNamespace(ctx, name)
		}()
	}

	wg.Wait()

	return nil
}
