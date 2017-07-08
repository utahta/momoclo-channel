package latestentry

import (
	"fmt"
	"strings"

	"github.com/utahta/momoclo-channel/model"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type repository struct{}

var Repository *repository = &repository{}

func (repo *repository) PutURL(ctx context.Context, url string) (*model.LatestEntry, error) {
	var code string
	blogCodes := []string{
		model.LatestEntryCodeTamai,
		model.LatestEntryCodeMomota,
		model.LatestEntryCodeAriyasu,
		model.LatestEntryCodeSasaki,
		model.LatestEntryCodeTakagi,
	}
	for _, c := range blogCodes {
		if strings.HasPrefix(url, fmt.Sprintf("https://ameblo.jp/%s", c)) {
			code = c
			break
		}
	}
	if strings.HasPrefix(url, "http://www.tfm.co.jp/clover/") {
		code = model.LatestEntryCodeHappyclo
	}

	if code == "" {
		// not supported
		return nil, nil
	}

	l := model.NewLatestEntry(code, "")
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
	l := model.NewLatestEntry(code, "")
	if err := l.Get(ctx); err != nil {
		return ""
	}
	return l.URL
}

func (repo *repository) GetTamaiURL(ctx context.Context) string {
	return repo.getURL(ctx, model.LatestEntryCodeTamai)
}

func (repo *repository) GetMomotaURL(ctx context.Context) string {
	return repo.getURL(ctx, model.LatestEntryCodeMomota)
}

func (repo *repository) GetAriyasuURL(ctx context.Context) string {
	return repo.getURL(ctx, model.LatestEntryCodeAriyasu)
}

func (repo *repository) GetSasakiURL(ctx context.Context) string {
	return repo.getURL(ctx, model.LatestEntryCodeSasaki)
}

func (repo *repository) GetTakagiURL(ctx context.Context) string {
	return repo.getURL(ctx, model.LatestEntryCodeTakagi)
}

func (repo *repository) GetHappycloURL(ctx context.Context) string {
	return repo.getURL(ctx, model.LatestEntryCodeHappyclo)
}
