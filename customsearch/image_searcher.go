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
	// ImageSearchResult represents custom search image result
	ImageSearchResult struct {
		URL          string
		ThumbnailURL string
	}

	// ImageSearcher represents google custom search image api
	ImageSearcher interface {
		Search(string) (ImageSearchResult, error)
	}

	imageSearcher struct {
		*customsearch.Service
	}
)

// NewImageSearcher returns CustomSearchClient
func NewImageSearcher(ctx context.Context) ImageSearcher {
	service, _ := customsearch.New(urlfetch.Client(ctx))
	return &imageSearcher{service}
}

// Search searches image given word
func (c *imageSearcher) Search(word string) (ImageSearchResult, error) {
	rand.Seed(time.Now().UnixNano())

	key := apiKey(config.C.GoogleCustomSearch.ApiKey)
	search, err := c.Cse.List(word).Cx(config.C.GoogleCustomSearch.ApiID).SearchType("image").Num(10).Start(rand.Int63n(30)).Do(key)
	if err != nil {
		return ImageSearchResult{}, err
	}

	for _, i := range rand.Perm(len(search.Items)) {
		item := search.Items[i]
		if item.Mime != "image/jpeg" {
			continue
		}
		if !strings.HasPrefix(item.Link, "https") {
			continue
		}
		return ImageSearchResult{URL: item.Link, ThumbnailURL: item.Image.ThumbnailLink}, nil
	}
	return ImageSearchResult{}, types.ErrNoSuchEntity
}
