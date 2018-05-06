package dao

import (
	"context"

	"github.com/mjibson/goon"
	"google.golang.org/appengine/datastore"
)

type (
	// Transactor provides transaction across entities
	Transactor interface {
		RunInTransaction(context.Context, func(context.Context) error, *TransactionOptions) error
	}

	// TransactionOptions represents transaction options (TODO: datastore dependencies should be eliminated, but there is no idea)
	TransactionOptions struct {
		XG       bool
		Attempts int
	}

	// datastoreTransactor implements Transactor interface using goon
	datastoreTransactor struct {
	}
)

// NewDatastoreTransactor wraps datastore transaction
func NewDatastoreTransactor() Transactor {
	return &datastoreTransactor{}
}

// RunInTransaction represents datastore transaction
func (t *datastoreTransactor) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error, opts *TransactionOptions) error {
	o := &datastore.TransactionOptions{XG: true}
	if opts != nil {
		o = &datastore.TransactionOptions{
			XG:       opts.XG,
			Attempts: opts.Attempts,
		}
	}

	return FromContext(ctx).RunInTransaction(func(g *goon.Goon) error {
		return fn(WithGoon(ctx, g))
	}, o)
}
