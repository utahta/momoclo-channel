package linenotification

import (
	"github.com/mjibson/goon"
	"github.com/utahta/momoclo-channel/model"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type repository struct{}

var Repository repository

func (repo repository) GetAll(ctx context.Context) ([]*model.LineNotification, error) {
	g := goon.FromContext(ctx)

	items := []*model.LineNotification{}
	q := datastore.NewQuery("LineNotification")

	_, err := g.GetAll(q, &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (repo repository) PutToken(ctx context.Context, token string) (*model.LineNotification, error) {
	ln, err := model.NewLineNotification(token)
	if err != nil {
		return nil, err
	}

	if err := ln.Put(ctx); err != nil {
		return nil, err
	}
	return ln, nil // save to datastore
}
