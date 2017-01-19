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
		ID:   code,
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

func (l *LatestBlogPost) Get(ctx context.Context) error {
	g := goon.FromContext(ctx)
	return g.Get(l)
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

	l := NewLatestBlogPost(code, "")
	if err := l.Get(ctx); err != nil && err != datastore.ErrNoSuchEntity {
		return nil, err
	}

	l.URL = url
	if err := l.Put(ctx); err != nil {
		return nil, err
	}
	return l, nil
}

func getLatestBlogPostURL(ctx context.Context, code string) string {
	q := datastore.NewQuery("LatestBlogPost").Filter("Code =", code)

	var posts []*LatestBlogPost
	if _, err := q.GetAll(ctx, &posts); err != nil {
		return ""
	}

	if len(posts) == 0 {
		return ""
	}
	return posts[0].URL
}

func GetTamaiLatestBlogPostURL(ctx context.Context) string {
	return getLatestBlogPostURL(ctx, BlogPostCodeTamai)
}

func GetMomotaLatestBlogPostURL(ctx context.Context) string {
	return getLatestBlogPostURL(ctx, BlogPostCodeMomota)
}

func GetAriyasuLatestBlogPostURL(ctx context.Context) string {
	return getLatestBlogPostURL(ctx, BlogPostCodeAriyasu)
}

func GetSasakiLatestBlogPostURL(ctx context.Context) string {
	return getLatestBlogPostURL(ctx, BlogPostCodeSasaki)
}

func GetTakagiLatestBlogPostURL(ctx context.Context) string {
	return getLatestBlogPostURL(ctx, BlogPostCodeTakagi)
}
