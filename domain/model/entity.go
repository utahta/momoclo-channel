package model

import (
	"time"

	"github.com/utahta/momoclo-channel/lib/timeutil"
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

	// PersistenceBeforeSaver hook
	PersistenceBeforeSaver interface {
		BeforeSave()
	}

	// PersistenceHandler represents persist operations
	PersistenceHandler interface {
		Put(interface{}) error
		PutMulti(interface{}) error
		Get(dst interface{}) error
		GetMulti(dst interface{}) error
		FlushLocalCache()
	}

	// TransactionOptions represents transaction options (TODO: want to eliminate datastore dependence but no ideas)
	TransactionOptions struct {
		XG       bool
		Attempts int
	}

	// Transactor provides transaction across entities
	Transactor interface {
		RunInTransaction(fn func(h PersistenceHandler) error, opts *TransactionOptions) error
		With(h PersistenceHandler, repos ...interface{}) (done func())
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
