package crawler

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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
	Url         string
	Title       string
	PublishedAt *time.Time
	Images      []*ChannelImage
	Videos      []*ChannelVideo
}

func (c *ChannelItem) UniqId() string {
	id := c.Url
	if c.PublishedAt != nil {
		id = fmt.Sprintf("%s%s", id, c.PublishedAt.Format("20060102150405"))
	}
	return id
}

func (c *ChannelItem) ImageUrlsToString() string {
	s := []string{}
	for _, image := range c.Images {
		s = append(s, image.Url)
	}
	return strings.Join(s, ",")
}

func (c *ChannelItem) VideoUrlsToString() string {
	s := []string{}
	for _, video := range c.Videos {
		s = append(s, video.Url)
	}
	return strings.Join(s, ",")
}

type Channel struct {
	Url    string
	Title  string
	Items  []*ChannelItem
	Client *http.Client `json:"-"`
}

type ChannelClient struct {
	Channel *Channel
	parser  ChannelParser
}

func newChannel(url string, title string) *Channel {
	return &Channel{
		Url:    url,
		Title:  title,
		Client: http.DefaultClient,
	}
}

func newChannelClient(c *Channel, parser ChannelParser) *ChannelClient {
	return &ChannelClient{
		Channel: c,
		parser:  parser,
	}
}

func (c *ChannelClient) Fetch() (*Channel, error) {
	resp, err := c.Channel.Client.Get(c.Channel.Url)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get resource. url:%s", c.Channel.Url)
	}
	defer resp.Body.Close()

	if c.parser == nil {
		return nil, errors.Errorf("You must implemented ChannelParser. url:%s", c.Channel.Url)
	}

	c.Channel.Items, err = c.parser.Parse(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse channel. url:%s", c.Channel.Url)
	}
	return c.Channel, nil
}
