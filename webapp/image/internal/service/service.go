package service

import (
	"context"
	"k8soperation/image/internal/models"
)

func ListImages(ctx context.Context) (any, error) {
	images, err := models.ListImages(ctx)
	if err != nil {
		return nil, err
	}

	return images, nil
}
