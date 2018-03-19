package customsearch

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/utahta/momoclo-channel/config"
	"github.com/utahta/momoclo-channel/domain"
	"github.com/utahta/momoclo-channel/domain/model"
	"google.golang.org/api/customsearch/v1"
	"google.golang.org/appengine/urlfetch"
)

type (
	imageSearcher struct {
		*customsearch.Service
	}

	apiKey string
)

// NewImageSearcher returns CustomSearchClient
func NewImageSearcher(ctx context.Context) model.ImageSearcher {
	service, _ := customsearch.New(urlfetch.Client(ctx))
	return &imageSearcher{service}
}

// Search searches image given word
func (c *imageSearcher) Search(word string) (model.ImageSearchResult, error) {
	rand.Seed(time.Now().UnixNano())

	key := apiKey(config.C.GoogleCustomSearch.ApiKey)
	search, err := c.Cse.List(word).Cx(config.C.GoogleCustomSearch.ApiID).SearchType("image").Num(10).Start(rand.Int63n(30)).Do(key)
	if err != nil {
		return model.ImageSearchResult{}, err
	}

	for _, i := range rand.Perm(len(search.Items)) {
		item := search.Items[i]
		if item.Mime != "image/jpeg" {
			continue
		}
		if !strings.HasPrefix(item.Link, "https") {
			continue
		}
		return model.ImageSearchResult{URL: item.Link, ThumbnailURL: item.Image.ThumbnailLink}, nil
	}
	return model.ImageSearchResult{}, domain.ErrNoSuchEntity
}

func (a apiKey) Get() (string, string) {
	return "key", string(a)
}
