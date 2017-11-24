package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type (
	// LatestEntry for confirm last updated entry url
	LatestEntry struct {
		ID        string    `datastore:"-" goon:"id"`
		Code      FeedCode  `validate:"required"`
		URL       string    `validate:"required,url"`
		CreatedAt time.Time `validate:"required"`
		UpdatedAt time.Time `validate:"required"`
	}

	// LatestEntryRepository interface
	LatestEntryRepository interface {
		Save(*LatestEntry) error
		FindOrNewByURL(string) (*LatestEntry, error)
		GetURL(string) string
	}
)

// NewLatestEntry builds *LatestEntry given url
func NewLatestEntry(urlStr string) (*LatestEntry, error) {
	var code FeedCode
	blogCodes := []FeedCode{
		FeedCodeTamai,
		FeedCodeMomota,
		FeedCodeAriyasu,
		FeedCodeSasaki,
		FeedCodeTakagi,
	}
	for _, c := range blogCodes {
		if strings.HasPrefix(urlStr, fmt.Sprintf("https://ameblo.jp/%s", c)) {
			code = c
			break
		}
	}
	if strings.HasPrefix(urlStr, "http://www.tfm.co.jp/clover/") {
		code = FeedCodeHappyclo
	} else if strings.HasPrefix(urlStr, "http://www.momoclo.net") {
		code = FeedCodeAeNews
	} else if strings.HasPrefix(urlStr, "https://www.youtube.com") {
		code = FeedCodeYoutube
	}

	if code == "" {
		return nil, errors.New("latestentry.Parse failed: code not supported")
	}
	return &LatestEntry{ID: string(code), Code: code, URL: urlStr}, nil
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
