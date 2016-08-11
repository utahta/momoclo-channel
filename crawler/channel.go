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

type Channel struct {
	Url string
	HttpClient *http.Client
}

type ChannelClient struct {
	Channel *Channel
	parser ChannelParser
}

func newChannel(url string) *Channel {
	return &Channel{
		Url: url,
		HttpClient: http.DefaultClient,
	}
}

func newChannelClient(c *Channel, parser ChannelParser) *ChannelClient {
	return &ChannelClient{
		Channel: c,
		parser: parser,
	}
}

func (c *ChannelClient) Fetch() ([]*ChannelItem, error) {
	resp, err := c.Channel.HttpClient.Get(c.Channel.Url)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get resource. url:%s", c.Channel.Url)
	}
	defer resp.Body.Close()

	if c.parser == nil {
		return nil, errors.Errorf("You must implemented ChannelParser. url:%s", c.Channel.Url)
	}
	return c.parser.Parse(resp.Body)
}
