package persistence

import "google.golang.org/appengine/datastore"

// DatastoreHandler interface
type DatastoreHandler interface {
	Put(interface{}) error
	PutMulti(interface{}) error
	Get(dst interface{}) error
	GetMulti(dst interface{}) error
	RunInTransaction(fn func(h DatastoreHandler) error, opts *datastore.TransactionOptions) error
}
