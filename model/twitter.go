package model

import (
	"time"

	"github.com/mjibson/goon"
	"github.com/utahta/momoclo-channel/crawler"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type TweetItem struct {
	Id          string `datastore:"-" goon:"id"`
	Title       string
	Url         string
	PublishedAt time.Time
	ImageUrls   string
	VideoUrls   string
}

func newTweetItem(item *crawler.ChannelItem) *TweetItem {
	return &TweetItem{
		Id:          item.UniqId(),
		Title:       item.Title,
		Url:         item.Url,
		PublishedAt: *item.PublishedAt,
		ImageUrls:   item.ImageUrlsToString(),
		VideoUrls:   item.VideoUrlsToString(),
	}
}

func PutTweetItem(ctx context.Context, item *crawler.ChannelItem) error {
	ti := newTweetItem(item)
	return datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		g := goon.FromContext(ctx)

		err := g.Get(ti)
		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}

		_, err = g.Put(ti)
		return err
	}, nil)
}
