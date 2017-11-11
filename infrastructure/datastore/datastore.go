package datastore

import (
	"context"

	"github.com/mjibson/goon"
	"github.com/utahta/momoclo-channel/adapter/persistence"
	"github.com/utahta/momoclo-channel/domain"
	"google.golang.org/appengine/datastore"
)

// handler implements DatastoreHandler interface using goon
type handler struct {
	*goon.Goon
}

// New returns DatastoreHandler
func New(c context.Context) persistence.DatastoreHandler {
	return &handler{
		goon.FromContext(c),
	}
}

// Put wraps goon.Put()
func (h *handler) Put(src interface{}) error {
	_, err := h.Goon.Put(src)
	return err
}

// PutMulti wraps goon.PutMulti()
func (h *handler) PutMulti(src interface{}) error {
	_, err := h.Goon.PutMulti(src)
	return err
}

// Get wraps goon.Get()
func (h *handler) Get(dst interface{}) error {
	err := h.Goon.Get(dst)
	if err == datastore.ErrNoSuchEntity {
		return domain.ErrNoSuchEntity
	}
	return err
}

// GetMulti wraps goon.GetMulti()
func (h *handler) GetMulti(dst interface{}) error {
	return h.Goon.GetMulti(dst)
}

// RunInTransaction wraps goon.RunInTransaction()
func (h *handler) RunInTransaction(fn func(h persistence.DatastoreHandler) error, opts *datastore.TransactionOptions) error {
	return h.Goon.RunInTransaction(func(g *goon.Goon) error {
		return fn(&handler{g})
	}, opts)
}
