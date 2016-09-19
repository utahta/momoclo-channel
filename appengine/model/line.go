package model

import (
	"time"

	"github.com/mjibson/goon"
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/crawler"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type LineItem struct {
	Id          string `datastore:"-" goon:"id"`
	Title       string
	Url         string
	PublishedAt time.Time
	ImageUrls   string `datastore:",noindex"`
	VideoUrls   string `datastore:",noindex"`
	CreatedAt   time.Time
}

func NewLineItem(item *crawler.ChannelItem) *LineItem {
	return &LineItem{
		Id:          item.UniqId(),
		Title:       item.Title,
		Url:         item.Url,
		PublishedAt: *item.PublishedAt,
		ImageUrls:   item.ImageUrlsToString(),
		VideoUrls:   item.VideoUrlsToString(),
	}
}

func (ti *LineItem) Put(ctx context.Context) error {
	g := goon.FromContext(ctx)

	// check for cached item
	err := g.Get(ti)
	if err == nil {
		return errors.Errorf("LineItem already exists.")
	} else if err != datastore.ErrNoSuchEntity {
		return err
	}

	return g.RunInTransaction(func(g *goon.Goon) error {
		err := g.Get(ti)
		if err != datastore.ErrNoSuchEntity {
			return err
		}

		jst, err := time.LoadLocation("Asia/Tokyo")
		if err != nil {
			return err
		}
		ti.CreatedAt = time.Now().In(jst)

		_, err = g.Put(ti)
		return err
	}, nil)
}
