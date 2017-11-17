package model

import (
	"time"

	"google.golang.org/appengine/datastore"
)

type (
	// TweetItem stores tweets
	TweetItem struct {
		ID          string `datastore:"-" goon:"id"`
		Title       string
		URL         string
		PublishedAt time.Time
		ImageURLs   string `datastore:",noindex"`
		VideoURLs   string `datastore:",noindex"`
		CreatedAt   time.Time
	}

	// TweetRequest represents tweet message, img and video data
	TweetRequest struct {
		InReplyToStatusID string
		Text              string
		ImageURLs         []string
		VideoURL          string
	}

	// TweetResponse represents tweet response data
	TweetResponse struct {
		IDStr string
	}

	// Tweeter interface
	Tweeter interface {
		Tweet(TweetRequest) (TweetResponse, error)
	}
)

// SetCreatedAt sets given time to CreatedAt
func (e *TweetItem) SetCreatedAt(t time.Time) {
	e.CreatedAt = t
}

// GetCreatedAt gets CreatedAt
func (e *TweetItem) GetCreatedAt() time.Time {
	return e.CreatedAt
}

// Load loads own from datastore
func (e *TweetItem) Load(p []datastore.Property) error {
	return load(e, p)
}

// Save saves own to datastore
func (e *TweetItem) Save() ([]datastore.Property, error) {
	return save(e)
}

// NewTweetItem returns TweetItem
//func NewTweetItem(item *crawler.ChannelItem) *TweetItem {
//	return &TweetItem{
//		ID:          item.UniqId(),
//		Title:       item.Title,
//		URL:         item.Url,
//		PublishedAt: *item.PublishedAt,
//		ImageURLs:   item.ImageUrlsToString(),
//		VideoURLs:   item.VideoUrlsToString(),
//	}
//}

// Put puts tweet item
//func (ti *TweetItem) Put(ctx context.Context) error {
//	g := goon.FromContext(ctx)
//
//	// check for cached item
//	err := g.Get(ti)
//	if err == nil {
//		return errors.Errorf("TweetItem already exists.")
//	} else if err != datastore.ErrNoSuchEntity {
//		return err
//	}
//
//	return g.RunInTransaction(func(g *goon.Goon) error {
//		err := g.Get(ti)
//		if err != datastore.ErrNoSuchEntity {
//			return err
//		}
//
//		jst, err := time.LoadLocation("Asia/Tokyo")
//		if err != nil {
//			return err
//		}
//		ti.CreatedAt = time.Now().In(jst)
//
//		_, err = g.Put(ti)
//		return err
//	}, nil)
//}
