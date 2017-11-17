package model

import (
	"time"
)

type (
	// FeedItem represents an entry in the feed
	FeedItem struct {
		Title       string
		URL         string
		EntryTitle  string
		EntryURL    string
		ImageURLs   []string
		VideoURLs   []string
		PublishedAt time.Time
	}

	// FeedFetcher interface
	FeedFetcher interface {
		Fetch(code string, maxItemNum int, latestURL string) ([]FeedItem, error)
	}
)
