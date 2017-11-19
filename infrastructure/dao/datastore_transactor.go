package dao

import (
	"context"

	"github.com/mjibson/goon"
	"github.com/utahta/momoclo-channel/domain/model"
	"google.golang.org/appengine/datastore"
)

type (
	// datastoreTransactor implements Transactor interface using goon
	datastoreTransactor struct {
		*goon.Goon
	}
)

// NewDatastoreTransactor wraps datastore transaction
func NewDatastoreTransactor(ctx context.Context) model.Transactor {
	return &datastoreTransactor{
		goon.FromContext(ctx),
	}
}

// RunInTransaction represents datastore transaction
func (t *datastoreTransactor) RunInTransaction(fn func(h model.PersistenceHandler) error, opts *model.TransactionOptions) error {
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

// With can be used in RunInTransaction
func (t *datastoreTransactor) With(h model.PersistenceHandler, args ...interface{}) (done func()) {
	dh, ok := h.(*datastoreHandler)
	if !ok {
		return func() {}
	}

	tmp := make([]*goon.Goon, len(args))
	for i, arg := range args {
		if v, ok := arg.(*datastoreHandler); ok {
			tmp[i] = v.Goon
			v.Goon = dh.Goon // deliver goon in transaction
		}
	}

	return func() {
		for i, arg := range args {
			if v, ok := arg.(*datastoreHandler); ok {
				v.Goon = tmp[i] // recovery
			}
		}
	}
}
