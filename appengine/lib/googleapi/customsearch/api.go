package customsearch

import (
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
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

	var key apiKey = apiKey(os.Getenv("GOOGLE_CUSTOM_SEARCH_API_KEY"))
	search, err := service.Cse.List(word).Cx(os.Getenv("GOOGLE_CUSTOM_SEARCH_API_ID")).SearchType("image").Num(10).Start(rand.Int63n(30)).Do(key)
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
