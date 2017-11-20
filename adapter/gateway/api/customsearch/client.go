package customsearch

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/utahta/momoclo-channel/domain"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/lib/config"
	"google.golang.org/api/customsearch/v1"
	"google.golang.org/appengine/urlfetch"
)

type (
	client struct {
		*customsearch.Service
	}

	apiKey string
)

// NewClient returns CustomSearchClient
func NewClient(ctx context.Context) model.CustomSearchClient {
	service, _ := customsearch.New(urlfetch.Client(ctx))
	return &client{service}
}

// SearchImage searches image given word
func (c *client) SearchImage(word string) (model.CustomSearchImageResult, error) {
	rand.Seed(time.Now().UnixNano())

	key := apiKey(config.C.GoogleCustomSearch.ApiKey)
	search, err := c.Cse.List(word).Cx(config.C.GoogleCustomSearch.ApiID).SearchType("image").Num(10).Start(rand.Int63n(30)).Do(key)
	if err != nil {
		return model.CustomSearchImageResult{}, err
	}

	for _, i := range rand.Perm(len(search.Items)) {
		item := search.Items[i]
		if item.Mime != "image/jpeg" {
			continue
		}
		if !strings.HasPrefix(item.Link, "https") {
			continue
		}
		return model.CustomSearchImageResult{URL: item.Link, ThumbnailURL: item.Image.ThumbnailLink}, nil
	}
	return model.CustomSearchImageResult{}, domain.ErrNoSuchEntity
}

func (a apiKey) Get() (string, string) {
	return "key", string(a)
}
