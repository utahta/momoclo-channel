package linenotification

import (
	"context"

	"github.com/mjibson/goon"
	"github.com/utahta/momoclo-channel/domain"
	"google.golang.org/appengine/datastore"
)

type repository struct{}

var Repository *repository = &repository{}

func (repo *repository) GetAll(ctx context.Context) ([]*domain.LineNotification, error) {
	g := goon.FromContext(ctx)

	items := []*domain.LineNotification{}
	q := datastore.NewQuery("LineNotification")

	_, err := g.GetAll(q, &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *repository) PutToken(ctx context.Context, token string) (*domain.LineNotification, error) {
	ln, err := domain.NewLineNotification(token)
	if err != nil {
		return nil, err
	}

	if err := ln.Put(ctx); err != nil {
		return nil, err
	}
	return ln, nil
}
