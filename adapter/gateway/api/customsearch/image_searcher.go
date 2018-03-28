package customsearch

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/utahta/momoclo-channel/config"
	"github.com/utahta/momoclo-channel/types"
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
func NewImageSearcher(ctx context.Context) types.ImageSearcher {
	service, _ := customsearch.New(urlfetch.Client(ctx))
	return &imageSearcher{service}
}

// Search searches image given word
func (c *imageSearcher) Search(word string) (types.ImageSearchResult, error) {
	rand.Seed(time.Now().UnixNano())

	key := apiKey(config.C.GoogleCustomSearch.ApiKey)
	search, err := c.Cse.List(word).Cx(config.C.GoogleCustomSearch.ApiID).SearchType("image").Num(10).Start(rand.Int63n(30)).Do(key)
	if err != nil {
		return types.ImageSearchResult{}, err
	}

	for _, i := range rand.Perm(len(search.Items)) {
		item := search.Items[i]
		if item.Mime != "image/jpeg" {
			continue
		}
		if !strings.HasPrefix(item.Link, "https") {
			continue
		}
		return types.ImageSearchResult{URL: item.Link, ThumbnailURL: item.Image.ThumbnailLink}, nil
	}
	return types.ImageSearchResult{}, types.ErrNoSuchEntity
}

func (a apiKey) Get() (string, string) {
	return "key", string(a)
}
