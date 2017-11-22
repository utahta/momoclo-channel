package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
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

// NewLatestEntry builds *LatestEntry given url
func NewLatestEntry(urlStr string) (*LatestEntry, error) {
	var code string
	blogCodes := []string{
		LatestEntryCodeTamai,
		LatestEntryCodeMomota,
		LatestEntryCodeAriyasu,
		LatestEntryCodeSasaki,
		LatestEntryCodeTakagi,
	}
	for _, c := range blogCodes {
		if strings.HasPrefix(urlStr, fmt.Sprintf("https://ameblo.jp/%s", c)) {
			code = c
			break
		}
	}
	if strings.HasPrefix(urlStr, "http://www.tfm.co.jp/clover/") {
		code = LatestEntryCodeHappyclo
	} else if strings.HasPrefix(urlStr, "http://www.momoclo.net") {
		code = LatestEntryCodeAeNews
	} else if strings.HasPrefix(urlStr, "https://www.youtube.com") {
		code = LatestEntryCodeYoutube
	}

	if code == "" {
		return nil, errors.New("latestentry.Parse failed: code not supported")
	}
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
