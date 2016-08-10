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

	Parse func(c *Channel, r io.Reader) ([]*ChannelItem, error)
}

func (c *Channel) Fetch() ([]*ChannelItem, error) {
	resp, err := c.HttpClient.Get(c.Url)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get resource. url:%s", c.Url)
	}
	defer resp.Body.Close()

	if c.Parse == nil {
		return nil, errors.Errorf("You must implemented Parse method. url:%s", c.Url)
	}

	return c.Parse(c, resp.Body)
}
