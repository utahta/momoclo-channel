package entity

import (
	"time"
)

type (
	// LatestEntry for confirm last updated entry url
	LatestEntry struct {
		ID          string    `datastore:"-" goon:"id"`
		Code        string    `validate:"required"`
		URL         string    `validate:"required,url"`
		PublishedAt time.Time `validate:"required"`
		CreatedAt   time.Time `validate:"required"`
		UpdatedAt   time.Time `validate:"required"`
	}
)

// NewLatestEntry builds *LatestEntry given url
func NewLatestEntry(code string, urlStr string) (*LatestEntry, error) {
	return &LatestEntry{ID: code, Code: code, URL: urlStr}, nil
}

// SetCreatedAt sets given time to CreatedAt
func (l *LatestEntry) SetCreatedAt(t time.Time) {
	l.CreatedAt = t
}

// GetCreatedAt gets CreatedAt
func (l *LatestEntry) GetCreatedAt() time.Time {
	return l.CreatedAt
}

// SetUpdatedAt sets given time to UpdatedAt
func (l *LatestEntry) SetUpdatedAt(t time.Time) {
	l.UpdatedAt = t
}

// BeforeSave hook
func (l *LatestEntry) BeforeSave() {
	beforeSave(l)
}
