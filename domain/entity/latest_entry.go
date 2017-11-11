package entity

import (
	"time"

	"google.golang.org/appengine/datastore"
)

const (
	LatestEntryCodeTamai    = "tamai-sd"
	LatestEntryCodeMomota   = "momota-sd"
	LatestEntryCodeAriyasu  = "ariyasu-sd"
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
		GetTamaiURL() string
		GetMomotaURL() string
		GetAriyasuURL() string
		GetSasakiURL() string
		GetTakagiURL() string
		GetHappycloURL() string
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

// Load loads own from datastore
func (l *LatestEntry) Load(p []datastore.Property) error {
	return load(l, p)
}

// Save saves own to datastore
func (l *LatestEntry) Save() ([]datastore.Property, error) {
	return save(l)
}
