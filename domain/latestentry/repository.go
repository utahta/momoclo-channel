package latestentry

import (
	"context"
	"fmt"
	"strings"

	"github.com/utahta/momoclo-channel/domain"
	"google.golang.org/appengine/datastore"
)

type repository struct{}

var Repository *repository = &repository{}

func (repo *repository) PutURL(ctx context.Context, url string) (*domain.LatestEntry, error) {
	var code string
	blogCodes := []string{
		domain.LatestEntryCodeTamai,
		domain.LatestEntryCodeMomota,
		domain.LatestEntryCodeAriyasu,
		domain.LatestEntryCodeSasaki,
		domain.LatestEntryCodeTakagi,
	}
	for _, c := range blogCodes {
		if strings.HasPrefix(url, fmt.Sprintf("https://ameblo.jp/%s", c)) {
			code = c
			break
		}
	}
	if strings.HasPrefix(url, "http://www.tfm.co.jp/clover/") {
		code = domain.LatestEntryCodeHappyclo
	}

	if code == "" {
		// not supported
		return nil, nil
	}

	l := domain.NewLatestEntry(code, "")
	if err := l.Get(ctx); err != nil && err != datastore.ErrNoSuchEntity {
		return nil, err
	}

	l.URL = url
	if err := l.Put(ctx); err != nil {
		return nil, err
	}
	return l, nil
}

func (repo *repository) getURL(ctx context.Context, code string) string {
	l := domain.NewLatestEntry(code, "")
	if err := l.Get(ctx); err != nil {
		return ""
	}
	return l.URL
}

func (repo *repository) GetTamaiURL(ctx context.Context) string {
	return repo.getURL(ctx, domain.LatestEntryCodeTamai)
}

func (repo *repository) GetMomotaURL(ctx context.Context) string {
	return repo.getURL(ctx, domain.LatestEntryCodeMomota)
}

func (repo *repository) GetAriyasuURL(ctx context.Context) string {
	return repo.getURL(ctx, domain.LatestEntryCodeAriyasu)
}

func (repo *repository) GetSasakiURL(ctx context.Context) string {
	return repo.getURL(ctx, domain.LatestEntryCodeSasaki)
}

func (repo *repository) GetTakagiURL(ctx context.Context) string {
	return repo.getURL(ctx, domain.LatestEntryCodeTakagi)
}

func (repo *repository) GetHappycloURL(ctx context.Context) string {
	return repo.getURL(ctx, domain.LatestEntryCodeHappyclo)
}
