package dao

import (
	"context"

	"github.com/mjibson/goon"
	"github.com/utahta/momoclo-channel/domain"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/infrastructure/dao/hook"
	"google.golang.org/appengine/datastore"
)

type (
	// datastoreHandler implements PersistenceHandler interface using goon
	datastoreHandler struct {
		*goon.Goon
	}

	// datastoreTransactor implements Transactor interface using goon
	datastoreTransactor struct {
		*goon.Goon
	}
)

// NewDatastoreHandler returns PersistenceHandler
func NewDatastoreHandler(ctx context.Context) model.PersistenceHandler {
	return &datastoreHandler{
		goon.FromContext(ctx),
	}
}

// Put wraps goon.Put()
func (h *datastoreHandler) Put(src interface{}) error {
	hook.BeforeSave(src)
	_, err := h.Goon.Put(src)
	return err
}

// PutMulti wraps goon.PutMulti()
func (h *datastoreHandler) PutMulti(src interface{}) error {
	hook.BeforeSaveMulti(src)
	_, err := h.Goon.PutMulti(src)
	return err
}

// Get wraps goon.Get()
func (h *datastoreHandler) Get(dst interface{}) error {
	err := h.Goon.Get(dst)
	if err == datastore.ErrNoSuchEntity {
		return domain.ErrNoSuchEntity
	}
	return err
}

// GetMulti wraps goon.GetMulti()
func (h *datastoreHandler) GetMulti(dst interface{}) error {
	return h.Goon.GetMulti(dst)
}

// FlushLocalCache clears local caches
func (h *datastoreHandler) FlushLocalCache() {
	h.Goon.FlushLocalCache()
}

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
