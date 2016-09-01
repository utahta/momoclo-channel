package model

import (
	"time"

	"github.com/mjibson/goon"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type LineUser struct {
	Id        string `datastore:"-" goon:"id"`
	Enabled   bool
	CreatedAt time.Time
}

func NewLineUser(id string) *LineUser {
	return &LineUser{Id: id, Enabled: true}
}

func (u *LineUser) Get(ctx context.Context) error {
	g := goon.FromContext(ctx)
	if err := g.Get(u); err != nil && err != datastore.ErrNoSuchEntity {
		return err
	}
	return nil
}

func (u *LineUser) Put(ctx context.Context) error {
	g := goon.FromContext(ctx)
	return g.RunInTransaction(func(g *goon.Goon) error {
		tmp := NewLineUser(u.Id)
		if err := g.Get(tmp); err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}

		if u.CreatedAt.IsZero() {
			jst, err := time.LoadLocation("Asia/Tokyo")
			if err != nil {
				return err
			}
			u.CreatedAt = time.Now().In(jst)
		}

		_, err := g.Put(u)
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

func (u *LineUserQuery) GetIds() ([]string, error) {
	q := datastore.NewQuery("LineUser").Filter("Enabled =", true).KeysOnly().Limit(u.Limit)
	if u.cursor.String() != "" {
		q = q.Start(u.cursor)
	}

	ids := []string{}
	t := q.Run(u.context)
	for {
		k, err := t.Next(nil)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		ids = append(ids, k.StringID())
	}

	var err error
	u.cursor, err = t.Cursor()
	if err != nil {
		return nil, err
	}
	return ids, nil
}
