package customsearch

import (
	"math/rand"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/api/customsearch/v1"
	"google.golang.org/appengine/urlfetch"
)

type ResultImage struct {
	Url          string
	ThumbnailUrl string
}

type apiKey string

func (a apiKey) Get() (string, string) { return "key", string(a) }

func SearchImage(ctx context.Context, word string) (*ResultImage, error) {
	service, err := customsearch.New(urlfetch.Client(ctx))
	if err != nil {
		return nil, err
	}

	var key apiKey = apiKey(os.Getenv("GOOGLE_CUSTOM_SEARCH_API_KEY"))
	search, err := service.Cse.List(word).Cx(os.Getenv("GOOGLE_CUSTOM_SEARCH_API_ID")).SearchType("image").Num(10).Start(rand.Int63n(20)).Do(key)
	if err != nil {
		return nil, err
	}

	for _, item := range search.Items {
		if item.Mime != "image/jpeg" {
			continue
		}
		res := &ResultImage{}
		res.Url = item.Link
		res.ThumbnailUrl = item.Image.ThumbnailLink
		return res, nil
	}
	return nil, errors.Errorf("Image not found. word:%s", word)
}
