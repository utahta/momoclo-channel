package entity

import (
	"time"

	"google.golang.org/appengine/datastore"
)

const (
	LatestEntryCodeTamai    = "tamai-sd"  // LatestEntryCodeTamai defines Shiori Tamai blog code
	LatestEntryCodeMomota   = "momota-sd" // LatestEntryCodeTamai defines Kanako Momota blog code
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

func (l *LatestEntry) SetCreatedAt(t time.Time) {
	l.CreatedAt = t
}

func (l *LatestEntry) GetCreatedAt() time.Time {
	return l.CreatedAt
}

func (l *LatestEntry) SetUpdatedAt(t time.Time) {
	l.UpdatedAt = t
}

func (l *LatestEntry) Load(p []datastore.Property) error {
	return load(l, p)
}

func (l *LatestEntry) Save() ([]datastore.Property, error) {
	return save(l)
}
