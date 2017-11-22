package model

import (
	"fmt"
	"time"
)

type (
	// LineItem represents line notification history
	LineItem struct {
		ID          string `datastore:"-" goon:"id"`
		Title       string
		URL         string
		PublishedAt time.Time
		ImageURLs   string `datastore:",noindex"`
		VideoURLs   string `datastore:",noindex"`
		CreatedAt   time.Time
	}

	// LineItemRepository interface
	LineItemRepository interface {
		Exists(string) bool
		Find(string) (*LineItem, error)
		Save(*LineItem) error
	}
)

// BuildID builds unique id
func (e *LineItem) BuildID() {
	id := e.URL
	if !e.PublishedAt.IsZero() {
		id = fmt.Sprintf("%s%s", id, e.PublishedAt.Format("20060102150405"))
	}
	e.ID = id
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
