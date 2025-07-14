package models

import (
	"context"
	"database/sql"
	"errors"
	"k8soperation/core"
)

func ListImages(ctx context.Context) ([]Image, error) {

	var images []Image
	err := core.Db.SelectContext(ctx, &images, "SELECT * FROM images")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make([]Image, 0), nil
		}
		return nil, err
	}

	return images, nil
}
