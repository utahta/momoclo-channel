package dao

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/dao/hook"
	"google.golang.org/appengine/datastore"
)

type (
	// PersistenceHandler represents persist operations
	PersistenceHandler interface {
		Kind(context.Context, interface{}) string
		Put(context.Context, interface{}) error
		PutMulti(context.Context, interface{}) error
		Get(context.Context, interface{}) error
		GetMulti(context.Context, interface{}) error
		Delete(context.Context, interface{}) error
		DeleteMulti(context.Context, interface{}) error
		NewQuery(string) PersistenceQuery
		GetAll(context.Context, PersistenceQuery, interface{}) error
		FlushLocalCache(context.Context)
	}

	// PersistenceQuery interface
	PersistenceQuery interface {
		Filter(string, interface{}) PersistenceQuery
	}

	// TransactionOptions represents transaction options (TODO: datastore dependencies should be eliminated, but there is no idea)
	TransactionOptions struct {
		XG       bool
		Attempts int
	}

	// datastoreHandler implements PersistenceHandler interface using goon
	datastoreHandler struct {
	}
)

// NewDatastoreHandler returns PersistenceHandler
func NewDatastoreHandler() PersistenceHandler {
	return &datastoreHandler{}
}

// Kind returns datastore kind given src
func (h *datastoreHandler) Kind(ctx context.Context, src interface{}) string {
	return FromContext(ctx).Kind(src)
}

// Put wraps goon.Put()
func (h *datastoreHandler) Put(ctx context.Context, src interface{}) error {
	hook.BeforeSave(src)
	if err := hook.Validate(src); err != nil {
		return err
	}
	_, err := FromContext(ctx).Put(src)
	return err
}

// PutMulti wraps goon.PutMulti()
func (h *datastoreHandler) PutMulti(ctx context.Context, src interface{}) error {
	v := reflect.Indirect(reflect.ValueOf(src))
	if v.Kind() != reflect.Slice {
		return errors.New("value must be a slice")
	}

	for i := 0; i < v.Len(); i++ {
		hook.BeforeSave(v.Index(i).Interface())

		//TODO perhaps should return multi error
		if err := hook.Validate(v.Index(i).Interface()); err != nil {
			return err
		}
	}
	_, err := FromContext(ctx).PutMulti(src)
	return err
}

// Get wraps goon.Get()
func (h *datastoreHandler) Get(ctx context.Context, dst interface{}) error {
	err := FromContext(ctx).Get(dst)
	if err == datastore.ErrNoSuchEntity {
		return ErrNoSuchEntity
	}
	return err
}

// GetMulti wraps goon.GetMulti()
func (h *datastoreHandler) GetMulti(ctx context.Context, dst interface{}) error {
	return FromContext(ctx).GetMulti(dst)
}

// Delete wraps goon.Delete()
func (h *datastoreHandler) Delete(ctx context.Context, src interface{}) error {
	g := FromContext(ctx)
	return g.Delete(g.Key(src))
}

// DeleteMulti wraps goon.DeleteMulti()
func (h *datastoreHandler) DeleteMulti(ctx context.Context, src interface{}) error {
	//TODO want to encapsulate logic that get datastore keys
	v := reflect.Indirect(reflect.ValueOf(src))
	if v.Kind() != reflect.Slice {
		return errors.New("datastore: value must be a slice or pointer-to-slice")
	}
	l := v.Len()

	g := FromContext(ctx)
	keys := make([]*datastore.Key, l)
	for i := 0; i < l; i++ {
		vi := v.Index(i)
		keys[i] = g.Key(vi.Interface())
	}
	return g.DeleteMulti(keys)
}

// Query returns PersistenceQuery that wraps datastore query
func (h *datastoreHandler) NewQuery(kind string) PersistenceQuery {
	return NewQuery(kind)
}

// GetAll runs the query and returns all matches entities
func (h *datastoreHandler) GetAll(ctx context.Context, q PersistenceQuery, dst interface{}) error {
	v, ok := q.(*datastoreQuery)
	if !ok {
		return errors.New("required datastoreQuery")
	}

	_, err := FromContext(ctx).GetAll(v.Query, dst)
	return err
}

// FlushLocalCache clears local caches
func (h *datastoreHandler) FlushLocalCache(ctx context.Context) {
	FromContext(ctx).FlushLocalCache()
}
