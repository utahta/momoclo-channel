package customsearch

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/lib/config"
	"google.golang.org/api/customsearch/v1"
	"google.golang.org/appengine/urlfetch"
)

type ResultImage struct {
	Url          string
	ThumbnailUrl string
}

var (
	ErrorImageNotFound = errors.Errorf("Image not found.")
)

type apiKey string

func (a apiKey) Get() (string, string) {
	return "key", string(a)
}

func SearchImage(ctx context.Context, word string) (*ResultImage, error) {
	service, err := customsearch.New(urlfetch.Client(ctx))
	if err != nil {
		return nil, err
	}
	rand.Seed(time.Now().UnixNano())

	var key apiKey = apiKey(config.C.GoogleCustomSearch.ApiKey)
	search, err := service.Cse.List(word).Cx(config.C.GoogleCustomSearch.ApiID).SearchType("image").Num(10).Start(rand.Int63n(30)).Do(key)
	if err != nil {
		return nil, err
	}

	for _, i := range rand.Perm(len(search.Items)) {
		item := search.Items[i]
		if item.Mime != "image/jpeg" {
			continue
		}
		if !strings.HasPrefix(item.Link, "https") {
			continue
		}

		res := &ResultImage{}
		res.Url = item.Link
		res.ThumbnailUrl = item.Image.ThumbnailLink
		return res, nil
	}
	return nil, ErrorImageNotFound
}
