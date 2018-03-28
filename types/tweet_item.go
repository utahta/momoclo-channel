package types

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

	// TweetItemRepository interface
	TweetItemRepository interface {
		Exists(string) bool
		Find(string) (*TweetItem, error)
		Save(*TweetItem) error
		Tx(PersistenceHandler) TweetItemRepository
	}
)

// NewTweetItem returns TweetItem given FeedItem
func NewTweetItem(params FeedItem) *TweetItem {
	return &TweetItem{
		ID:          params.UniqueURL(),
		Title:       params.EntryTitle,
		URL:         params.EntryURL,
		PublishedAt: params.PublishedAt,
		ImageURLs:   strings.Join(params.ImageURLs, ","),
		VideoURLs:   strings.Join(params.VideoURLs, ","),
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
