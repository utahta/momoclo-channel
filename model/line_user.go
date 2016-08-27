package model

import (
	"time"

	"github.com/mjibson/goon"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type LineUser struct {
	Id        string `datastore:"-" goon:"id"`
	CreatedAt time.Time
}

func NewLineUser(id string) *LineUser {
	return &LineUser{Id: id}
}

func (u *LineUser) Put(ctx context.Context) error {
	g := goon.FromContext(ctx)

	// check for cached item
	if g.Get(u) == nil {
		return errors.Errorf("LineUser already exists.")
	}

	return g.RunInTransaction(func(g *goon.Goon) error {
		err := g.Get(u)
		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}

		jst, err := time.LoadLocation("Asia/Tokyo")
		if err != nil {
			return err
		}
		u.CreatedAt = time.Now().In(jst)

		_, err = g.Put(u)
		return err
	}, nil)
}

type LineUserQuery struct {
	context context.Context
	cursor  datastore.Cursor
	Limit   int
}

func NewLineUserQuery(ctx context.Context) *LineUserQuery {
	return &LineUserQuery{context: ctx, Limit: 100}
}

func (u *LineUserQuery) GetIds(cursor datastore.Cursor) ([]string, datastore.Cursor, error) {
	q := datastore.NewQuery("LineUser").KeysOnly().Limit(u.Limit)
	if cursor.String() != "" {
		q = q.Start(cursor)
	}

	ids := []string{}
	t := q.Run(u.context)
	for {
		k, err := t.Next(nil)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, cursor, err
		}
		ids = append(ids, k.StringID())
	}

	cursor, err := t.Cursor()
	if err != nil {
		return nil, cursor, err
	}
	return ids, cursor, nil
}
