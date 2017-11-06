package persistence

// DatastoreHandler
type DatastoreHandler interface {
	Put(interface{}) error
	PutMulti(interface{}) error
	Get(dst interface{}) error
	GetMulti(dst interface{}) error
	RunInTransaction(fn func(h DatastoreHandler) error, opts *TransactionOptions) error
}

// TransactionOptions refs appengine/datastore.TransactionOptions
type TransactionOptions struct {
	XG       bool
	Attempts int
}
