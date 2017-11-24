package model

import (
	"time"
)

const (
	FeedCodeMomota   FeedCode = "momota-sd"
	FeedCodeAriyasu  FeedCode = "ariyasu-sd"
	FeedCodeTamai    FeedCode = "tamai-sd"
	FeedCodeSasaki   FeedCode = "sasaki-sd"
	FeedCodeTakagi   FeedCode = "takagi-sd"
	FeedCodeHappyclo FeedCode = "happyclo"
	FeedCodeAeNews   FeedCode = "aenews"
	FeedCodeYoutube  FeedCode = "youtube"
)

type (
	// FeedCode represents identify code of feed
	FeedCode string

	// FeedItem represents an entry in the feed
	FeedItem struct {
		Title       string `validate:"required"`
		URL         string `validate:"required,url"`
		EntryTitle  string `validate:"required"`
		EntryURL    string `validate:"required,url"`
		ImageURLs   []string
		VideoURLs   []string
		PublishedAt time.Time `validate:"required"`
	}

	// FeedFetcher interface
	FeedFetcher interface {
		Fetch(code FeedCode, maxItemNum int, latestURL string) ([]FeedItem, error)
	}
)

// String returns string representation of FeedCode
func (f FeedCode) String() string {
	return string(f)
}
