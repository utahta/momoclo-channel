package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/mjibson/goon"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
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

func PutLatestEntry(ctx context.Context, url string) (*LatestEntry, error) {
	var code string
	blogCodes := []string{
		LatestEntryCodeTamai,
		LatestEntryCodeMomota,
		LatestEntryCodeAriyasu,
		LatestEntryCodeSasaki,
		LatestEntryCodeTakagi,
	}
	for _, c := range blogCodes {
		if strings.HasPrefix(url, fmt.Sprintf("http://ameblo.jp/%s", c)) {
			code = c
			break
		}
	}
	if strings.HasPrefix(url, "http://www.tfm.co.jp/clover/") {
		code = LatestEntryCodeHappyclo
	}

	if code == "" {
		// maybe not blog post
		return nil, nil
	}

	l := NewLatestEntry(code, "")
	if err := l.Get(ctx); err != nil && err != datastore.ErrNoSuchEntity {
		return nil, err
	}

	l.URL = url
	if err := l.Put(ctx); err != nil {
		return nil, err
	}
	return l, nil
}

func getLatestEntryURL(ctx context.Context, code string) string {
	l := NewLatestEntry(code, "")
	if err := l.Get(ctx); err != nil {
		return ""
	}
	return l.URL
}

func GetTamaiLatestEntryURL(ctx context.Context) string {
	return getLatestEntryURL(ctx, LatestEntryCodeTamai)
}

func GetMomotaLatestEntryURL(ctx context.Context) string {
	return getLatestEntryURL(ctx, LatestEntryCodeMomota)
}

func GetAriyasuLatestEntryURL(ctx context.Context) string {
	return getLatestEntryURL(ctx, LatestEntryCodeAriyasu)
}

func GetSasakiLatestEntryURL(ctx context.Context) string {
	return getLatestEntryURL(ctx, LatestEntryCodeSasaki)
}

func GetTakagiLatestEntryURL(ctx context.Context) string {
	return getLatestEntryURL(ctx, LatestEntryCodeTakagi)
}

func GetHappycloLatestEntryURL(ctx context.Context) string {
	return getLatestEntryURL(ctx, LatestEntryCodeHappyclo)
}
