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

type ChannelParser interface {
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

type ChannelContext struct {
	Url string
	HttpClient *http.Client
}

type Channel struct {
	Context *ChannelContext
	parser ChannelParser
}

func newChannelContext(url string) *ChannelContext {
	return &ChannelContext{
		Url: url,
		HttpClient: http.DefaultClient,
	}
}

func newChannel(ctx *ChannelContext, parser ChannelParser) *Channel {
	return &Channel{ Context: ctx, parser: parser }
}

func (c *Channel) Fetch() ([]*ChannelItem, error) {
	resp, err := c.Context.HttpClient.Get(c.Context.Url)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get resource. url:%s", c.Context.Url)
	}
	defer resp.Body.Close()

	if c.parser == nil {
		return nil, errors.Errorf("You must implemented ChannelParser. url:%s", c.Context.Url)
	}
	return c.parser.Parse(resp.Body)
}
