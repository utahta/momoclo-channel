package entity

import (
	"github.com/utahta/momoclo-channel/timeutil"
	"github.com/utahta/momoclo-channel/types"
)

func beforeSave(src interface{}) {
	now := timeutil.Now()

	if v, ok := src.(types.CreateTimestamper); ok {
		if v.GetCreatedAt().IsZero() {
			v.SetCreatedAt(now)
		}
	}

	if v, ok := src.(types.UpdateTimestamper); ok {
		v.SetUpdatedAt(now)
	}
}
