package model

import (
	"time"

	"github.com/utahta/momoclo-channel/lib/config"
	"google.golang.org/appengine/datastore"
)

// CreateTimestamper provides set current timestamp to entity when save if not already set
type CreateTimestamper interface {
	SetCreatedAt(time.Time)
	GetCreatedAt() time.Time
}

// UpdateTimestamper provides set current timestamp to entity when save
type UpdateTimestamper interface {
	SetUpdatedAt(time.Time)
}

func load(dst interface{}, p []datastore.Property) error {
	return datastore.LoadStruct(dst, p)
}

func save(src interface{}) ([]datastore.Property, error) {
	now := time.Now().In(config.JST)

	if v, ok := src.(CreateTimestamper); ok {
		if v.GetCreatedAt().IsZero() {
			v.SetCreatedAt(now)
		}
	}

	if v, ok := src.(UpdateTimestamper); ok {
		v.SetUpdatedAt(now)
	}

	return datastore.SaveStruct(src)
}
