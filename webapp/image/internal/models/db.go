package models

import (
	"context"
	"k8soperation/core"
)

func ListImages(ctx context.Context) ([]Image, error) {

	var images []Image
	err := core.Db.SelectContext(ctx, &images, "SELET * FROM images")
	if err != nil {
		return nil, err
	}

	return images, nil
}
