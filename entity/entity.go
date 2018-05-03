package entity

import (
	"time"

	"github.com/utahta/momoclo-channel/timeutil"
)

type (
	// CreateTimestamper provides set current timestamp to entity when save if not already set
	CreateTimestamper interface {
		SetCreatedAt(time.Time)
		GetCreatedAt() time.Time
	}

	// UpdateTimestamper provides set current timestamp to entity when save
	UpdateTimestamper interface {
		SetUpdatedAt(time.Time)
	}
)

func beforeSave(src interface{}) {
	now := timeutil.Now()

	if v, ok := src.(CreateTimestamper); ok {
		if v.GetCreatedAt().IsZero() {
			v.SetCreatedAt(now)
		}
	}

	if v, ok := src.(UpdateTimestamper); ok {
		v.SetUpdatedAt(now)
	}
}
