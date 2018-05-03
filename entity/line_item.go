package entity

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
)

// NewLineItem returns LineItem given FeedItem
func NewLineItem(id string, title string, url string, publishedAt time.Time, imageURLs []string, videoURLs []string) *LineItem {
	return &LineItem{
		ID:          id,
		Title:       title,
		URL:         url,
		PublishedAt: publishedAt,
		ImageURLs:   strings.Join(imageURLs, ","),
		VideoURLs:   strings.Join(videoURLs, ","),
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
