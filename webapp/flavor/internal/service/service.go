package service

import (
	"context"
	"k8soperation/flavor/internal/models"
)

func ListFlavors(ctx context.Context) (any, error) {
	flavors, err := models.ListFlavors(ctx)
	if err != nil {
		return nil, err
	}

	return flavors, nil
}
