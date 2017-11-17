package model

import (
	"time"
)

const (
	LatestEntryCodeMomota   = "momota-sd"
	LatestEntryCodeAriyasu  = "ariyasu-sd"
	LatestEntryCodeTamai    = "tamai-sd"
	LatestEntryCodeSasaki   = "sasaki-sd"
	LatestEntryCodeTakagi   = "takagi-sd"
	LatestEntryCodeHappyclo = "happyclo"
	LatestEntryCodeAeNews   = "aenews"
	LatestEntryCodeYoutube  = "youtube"
)

type (
	// LatestEntry for confirm last updated entry url
	LatestEntry struct {
		ID        string `datastore:"-" goon:"id"`
		Code      string
		URL       string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	// LatestEntryRepository interface
	LatestEntryRepository interface {
		Save(*LatestEntry) error
		FindOrCreateByURL(string) (*LatestEntry, error)
		GetURL(string) string
	}
)

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
