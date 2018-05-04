package entity

import (
	"strings"
	"time"
)

type (
	// TweetItem represents tweet history
	TweetItem struct {
		ID          string    `datastore:"-" goon:"id" validate:"required"`
		Title       string    `validate:"required"`
		URL         string    `validate:"required,url"`
		PublishedAt time.Time `validate:"required"`
		ImageURLs   string    `datastore:",noindex"`
		VideoURLs   string    `datastore:",noindex"`
		CreatedAt   time.Time `validate:"required"`
	}
)

// NewTweetItem returns TweetItem given FeedItem
func NewTweetItem(id string, title string, url string, publishedAt time.Time, imageURLs []string, videoURLs []string) *TweetItem {
	return &TweetItem{
		ID:          id,
		Title:       title,
		URL:         url,
		PublishedAt: publishedAt,
		ImageURLs:   strings.Join(imageURLs, ","),
		VideoURLs:   strings.Join(videoURLs, ","),
	}
}

// SetCreatedAt sets given time to CreatedAt
func (e *TweetItem) SetCreatedAt(t time.Time) {
	e.CreatedAt = t
}

// GetCreatedAt gets CreatedAt
func (e *TweetItem) GetCreatedAt() time.Time {
	return e.CreatedAt
}

// BeforeSave hook
func (e *TweetItem) BeforeSave() {
	beforeSave(e)
}
