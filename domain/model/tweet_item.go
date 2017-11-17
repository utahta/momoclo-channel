package model

import (
	"fmt"
	"time"
)

type (
	// TweetItem represents tweet history
	TweetItem struct {
		ID          string `datastore:"-" goon:"id"`
		Title       string
		URL         string
		PublishedAt time.Time
		ImageURLs   string `datastore:",noindex"`
		VideoURLs   string `datastore:",noindex"`
		CreatedAt   time.Time
	}

	// TweetItemRepository interface
	TweetItemRepository interface {
		Exists(string) bool
		Find(string) (*TweetItem, error)
		Save(*TweetItem) error
	}
)

// BuildID builds unique id
func (e *TweetItem) BuildID() {
	id := e.URL
	if !e.PublishedAt.IsZero() {
		id = fmt.Sprintf("%s%s", id, e.PublishedAt.Format("20060102150405"))
	}
	e.ID = id
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
