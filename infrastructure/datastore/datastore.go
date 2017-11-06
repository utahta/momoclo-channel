package datastore

import (
	"context"

	"github.com/mjibson/goon"
	"github.com/utahta/momoclo-channel/adapter/persistence"
	"github.com/utahta/momoclo-channel/domain"
	"google.golang.org/appengine/datastore"
)

// Handler implements datastoreHandler interface using goon
type Handler struct {
	*goon.Goon
}

// New returns datastore Handler
func New(c context.Context) persistence.DatastoreHandler {
	return &Handler{
		goon.FromContext(c),
	}
}

// Put wraps goon.Put()
func (h *Handler) Put(src interface{}) error {
	_, err := h.Goon.Put(src)
	return err
}

// PutMulti wraps goon.PutMulti()
func (h *Handler) PutMulti(src interface{}) error {
	_, err := h.Goon.PutMulti(src)
	return err
}

// Get wraps goon.Get()
func (h *Handler) Get(dst interface{}) error {
	err := h.Goon.Get(dst)
	if err == datastore.ErrNoSuchEntity {
		return domain.ErrNoSuchEntity
	}
	return err
}

// GetMulti wraps goon.GetMulti()
func (h *Handler) GetMulti(dst interface{}) error {
	return h.Goon.GetMulti(dst)
}

// RunInTransaction wraps goon.RunInTransaction()
func (h *Handler) RunInTransaction(fn func(h persistence.DatastoreHandler) error, o *persistence.TransactionOptions) error {
	var opts *datastore.TransactionOptions
	if o != nil {
		opts = &datastore.TransactionOptions{XG: o.XG, Attempts: o.Attempts}
	}

	return h.Goon.RunInTransaction(func(g *goon.Goon) error {
		return fn(&Handler{g})
	}, opts)
}
