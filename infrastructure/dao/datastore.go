package dao

import (
	"context"
	"reflect"

	"github.com/mjibson/goon"
	"github.com/pkg/errors"
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
)

// NewDatastoreHandler returns PersistenceHandler
func NewDatastoreHandler(ctx context.Context) model.PersistenceHandler {
	return &datastoreHandler{
		goon.FromContext(ctx),
	}
}

// Kind returns datastore kind given src
func (h *datastoreHandler) Kind(src interface{}) string {
	return h.Goon.Kind(src)
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

// Delete wraps goon.Delete()
func (h *datastoreHandler) Delete(src interface{}) error {
	return h.Goon.Delete(h.Goon.Key(src))
}

// DeleteMulti wraps goon.DeleteMulti()
func (h *datastoreHandler) DeleteMulti(src interface{}) error {
	//TODO want to encapsulate logic that get datastore keys
	v := reflect.Indirect(reflect.ValueOf(src))
	if v.Kind() != reflect.Slice {
		return errors.New("datastore: value must be a slice or pointer-to-slice")
	}
	l := v.Len()

	keys := make([]*datastore.Key, l)
	for i := 0; i < l; i++ {
		vi := v.Index(i)
		keys[i] = h.Goon.Key(vi.Interface())
	}
	return h.Goon.DeleteMulti(keys)
}

// Query returns PersistenceQuery that wraps datastore query
func (h *datastoreHandler) NewQuery(kind string) model.PersistenceQuery {
	return NewQuery(kind)
}

// GetAll runs the query and returns all matches entities
func (h *datastoreHandler) GetAll(q model.PersistenceQuery, dst interface{}) error {
	v, ok := q.(*datastoreQuery)
	if !ok {
		return errors.New("required datastoreQuery")
	}

	_, err := h.Goon.GetAll(v.Query, dst)
	return err
}

// FlushLocalCache clears local caches
func (h *datastoreHandler) FlushLocalCache() {
	h.Goon.FlushLocalCache()
}
