package domain

import (
	"context"
	"time"

	"github.com/mjibson/goon"
)

const (
	LatestEntryCodeTamai    = "tamai-sd"
	LatestEntryCodeMomota   = "momota-sd"
	LatestEntryCodeAriyasu  = "ariyasu-sd"
	LatestEntryCodeSasaki   = "sasaki-sd"
	LatestEntryCodeTakagi   = "takagi-sd"
	LatestEntryCodeHappyclo = "happyclo"
)

type LatestEntry struct {
	ID        string `datastore:"-" goon:"id"`
	Code      string
	URL       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewLatestEntry(code string, url string) *LatestEntry {
	return &LatestEntry{
		ID:   code,
		Code: code,
		URL:  url,
	}
}

func (l *LatestEntry) Put(ctx context.Context) error {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}
	now := time.Now().In(jst)

	g := goon.FromContext(ctx)
	return g.RunInTransaction(func(g *goon.Goon) error {
		if l.CreatedAt.IsZero() {
			l.CreatedAt = now
		}
		l.UpdatedAt = now

		_, err = g.Put(l)
		return err
	}, nil)
}

func (l *LatestEntry) Get(ctx context.Context) error {
	g := goon.FromContext(ctx)
	return g.Get(l)
}
