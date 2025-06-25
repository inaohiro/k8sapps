package models

import (
	"context"
	"k8soperation/core"
)

func ListFlavors(ctx context.Context) ([]Flavor, error) {

	var flavors []Flavor
	err := core.Db.SelectContext(ctx, &flavors, "SELET * FROM flavors")
	if err != nil {
		return nil, err
	}

	return flavors, nil
}

