package models

import (
	"context"
	"database/sql"
	"errors"
	"k8soperation/core"
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
	return result, nil
}
