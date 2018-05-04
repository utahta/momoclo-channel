package dao

import (
	"context"

	"github.com/mjibson/goon"
	"google.golang.org/appengine/datastore"
)

type (
	// datastoreTransactor implements Transactor interface using goon
	datastoreTransactor struct {
		*goon.Goon
	}
)

// NewDatastoreTransactor wraps datastore transaction
func NewDatastoreTransactor(ctx context.Context) Transactor {
	return &datastoreTransactor{
		goon.FromContext(ctx),
	}
}

// RunInTransaction represents datastore transaction
func (t *datastoreTransactor) RunInTransaction(fn func(h PersistenceHandler) error, opts *TransactionOptions) error {
	o := &datastore.TransactionOptions{XG: true}
	if opts != nil {
		o = &datastore.TransactionOptions{
			XG:       opts.XG,
			Attempts: opts.Attempts,
		}
	}

	return t.Goon.RunInTransaction(func(g *goon.Goon) error {
		return fn(&datastoreHandler{g})
	}, o)
}
