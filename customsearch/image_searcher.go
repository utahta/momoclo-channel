package customsearch

import (
	"context"
	"math/rand"
	"strings"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/config"
	"google.golang.org/api/customsearch/v1"
	"google.golang.org/appengine/urlfetch"
)

type (
	// ImageSearchResult represents the result of image search
	ImageSearchResult struct {
		URL          string
		ThumbnailURL string
	}

	// ImageSearcher represents image search that uses google custom search api
	ImageSearcher interface {
		Search(context.Context, string) (ImageSearchResult, error)
	}

	imageSearcher struct {
	}
)

// NewImageSearcher returns CustomSearchClient
func NewImageSearcher() ImageSearcher {
	return &imageSearcher{}
}

// Search searches image given word
func (c *imageSearcher) Search(ctx context.Context, word string) (ImageSearchResult, error) {
	s, err := customsearch.New(urlfetch.Client(ctx))
	if err != nil {
		return ImageSearchResult{}, err
	}

	key := apiKey(config.C().GoogleCustomSearch.ApiKey)
	search, err := s.Cse.List(word).Cx(config.C().GoogleCustomSearch.ApiID).SearchType("image").Num(10).Start(rand.Int63n(30)).Do(key)
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
	return ImageSearchResult{}, errors.New("image not found")
}
