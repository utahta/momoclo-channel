package dao

import (
	"context"

	"github.com/mjibson/goon"
	"github.com/utahta/momoclo-channel/adapter/gateway/persistence"
	"github.com/utahta/momoclo-channel/domain"
	"google.golang.org/appengine/datastore"
)

// datastoreHandler implements DatastoreHandler interface using goon
type datastoreHandler struct {
	*goon.Goon
}

// NewDatastoreHandler returns DatastoreHandler
func NewDatastoreHandler(c context.Context) persistence.DatastoreHandler {
	return &datastoreHandler{
		goon.FromContext(c),
	}
}

// Put wraps goon.Put()
func (h *datastoreHandler) Put(src interface{}) error {
	_, err := h.Goon.Put(src)
	return err
}

// PutMulti wraps goon.PutMulti()
func (h *datastoreHandler) PutMulti(src interface{}) error {
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

// RunInTransaction wraps goon.RunInTransaction()
func (h *datastoreHandler) RunInTransaction(fn func(h persistence.DatastoreHandler) error, opts *datastore.TransactionOptions) error {
	return h.Goon.RunInTransaction(func(g *goon.Goon) error {
		return fn(&datastoreHandler{g})
	}, opts)
}
