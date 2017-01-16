package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/mjibson/goon"
	"golang.org/x/net/context"
)

const (
	BlogPostCodeTamai   = "tamai-sd"
	BlogPostCodeMomota  = "momota-sd"
	BlogPostCodeAriyasu = "ariyasu-sd"
	BlogPostCodeSasaki  = "sasaki-sd"
	BlogPostCodeTakagi  = "takagi-sd"
)

type LatestBlogPost struct {
	ID        string `datastore:"-" goon:"id"`
	Code      string
	URL       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewLatestBlogPost(code string, url string) *LatestBlogPost {
	return &LatestBlogPost{
		ID:   fmt.Sprintf("%s-%s", code, url),
		Code: code,
		URL:  url,
	}
}

func (l *LatestBlogPost) Put(ctx context.Context) error {
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

func PutLatestBlogPost(ctx context.Context, url string) (*LatestBlogPost, error) {
	var code string
	codes := []string{
		BlogPostCodeTamai,
		BlogPostCodeMomota,
		BlogPostCodeAriyasu,
		BlogPostCodeSasaki,
		BlogPostCodeTakagi,
	}
	for _, c := range codes {
		if strings.HasPrefix(url, fmt.Sprintf("http://ameblo.jp/%s", c)) {
			code = c
			break
		}
	}

	if code == "" {
		// maybe not blog post
		return nil, nil
	}

	l := NewLatestBlogPost(code, url)
	if err := l.Put(ctx); err != nil {
		return nil, err
	}
	return l, nil
}
