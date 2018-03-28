package types

import (
	"strings"
	"time"
)

type (
	// LineItem represents line notification history
	LineItem struct {
		ID          string    `datastore:"-" goon:"id" validate:"required"`
		Title       string    `validate:"required"`
		URL         string    `validate:"required,url"`
		PublishedAt time.Time `validate:"required"`
		ImageURLs   string    `datastore:",noindex"`
		VideoURLs   string    `datastore:",noindex"`
		CreatedAt   time.Time `validate:"required"`
	}

	// LineItemRepository interface
	LineItemRepository interface {
		Exists(string) bool
		Find(string) (*LineItem, error)
		Save(*LineItem) error
		Tx(PersistenceHandler) LineItemRepository
	}
)

// NewLineItem returns LineItem given FeedItem
func NewLineItem(params FeedItem) *LineItem {
	return &LineItem{
		ID:          params.UniqueURL(),
		Title:       params.EntryTitle,
		URL:         params.EntryURL,
		PublishedAt: params.PublishedAt,
		ImageURLs:   strings.Join(params.ImageURLs, ","),
		VideoURLs:   strings.Join(params.VideoURLs, ","),
	}
}

// SetCreatedAt sets given time to CreatedAt
func (e *LineItem) SetCreatedAt(t time.Time) {
	e.CreatedAt = t
}

// GetCreatedAt gets CreatedAt
func (e *LineItem) GetCreatedAt() time.Time {
	return e.CreatedAt
}

// BeforeSave hook
func (e *LineItem) BeforeSave() {
	beforeSave(e)
}
