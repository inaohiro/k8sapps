package models

import (
	"context"
	"database/sql"
	"errors"
	"k8soperation/core"
	"math/rand/v2"
	"time"
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

	result := make([]Image, 0, len(images))
	for _, v := range images {
		result = append(result, Image{
			Name: v.Name,
		})
	}

	// DB の実行に 300ms ~ 500ms の遅延を発生させる
	time.Sleep(time.Duration(rand.IntN(200)+300) * time.Millisecond)

	return result, nil
}
