package models

import (
	"context"
	"database/sql"
	"errors"
	"k8soperation/core"
	"math/rand/v2"
	"time"
)

func ListFlavors(ctx context.Context) ([]Flavor, error) {
	var flavors []Flavor
	err := core.Db.SelectContext(ctx, &flavors, "SELECT * FROM flavors")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make([]Flavor, 0), nil
		}
		return nil, err
	}

	result := make([]Flavor, 0, len(flavors))
	for _, v := range flavors {
		result = append(result, Flavor{
			Name: v.Name,
		})
	}

	// DB の実行に 300ms ~ 500ms の遅延を発生させる
	time.Sleep(time.Duration(rand.IntN(200)+300) * time.Millisecond)

	return result, nil
}
