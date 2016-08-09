package crawler

import (
	"time"
	"net/http"
	"io"

	"github.com/pkg/errors"
)

const (
	timeout = 5 // sec
)

type ChannelFetchParser interface {
	Fetch() (io.ReadCloser, error)
	Parse(r io.Reader) ([]*ChannelItem, error)
}

type ChannelImage struct {
	Url string
}

type ChannelVideo struct {
	Url string
}

type ChannelItem struct {
	Url string
	Title string
	PublishedAt *time.Time
	Images []*ChannelImage
	Videos []*ChannelVideo
}

type Channel struct {
	Url string
	HttpClient http.Client
}

func (c *Channel) fetch(url string) (io.ReadCloser, error) {
	res, err := c.HttpClient.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to fetch url:%s", url)
	}
	return res.Body, nil
}

func FetchParse(c ChannelFetchParser) ([]*ChannelItem, error) {
	r, err := c.Fetch()
	if err != nil {
		return nil, err
	}
	return c.Parse(r)
}
