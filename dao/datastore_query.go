package dao

import (
	"github.com/utahta/momoclo-channel/types"
	"google.golang.org/appengine/datastore"
)

type (
	datastoreQuery struct {
		*datastore.Query
	}
)

// NewQuery returns PersistenceQuery wraps datastore.Query
func NewQuery(kind string) types.PersistenceQuery {
	return &datastoreQuery{
		Query: datastore.NewQuery(kind),
	}
}

// Filter wraps datastore.Query.Filter
func (q *datastoreQuery) Filter(filterStr string, value interface{}) types.PersistenceQuery {
	q.Query = q.Query.Filter(filterStr, value)
	return q
}
