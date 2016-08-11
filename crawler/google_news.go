package crawler

import (
	"time"
	"io"
	"io/ioutil"

	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/kennygrant/sanitize"
	"github.com/pkg/errors"
)

type googleNewsChannelParser struct {
	channel *Channel
}

func NewGoogleNewsChannelClient() *ChannelClient {
	c := newChannel("https://www.google.com/alerts/feeds/15513821572968738743/9316362605522861420")
	return newChannelClient(c, &googleNewsChannelParser{ channel: c })
}

func FetchGoogleNews() ([]*ChannelItem, error) {
	return NewGoogleNewsChannelClient().Fetch()
}

func (p *googleNewsChannelParser) Parse(r io.Reader) ([]*ChannelItem, error) {
	c := p.channel
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read rss content")
	}

	feed := rss.New(timeout, true, nil, nil)
	err = feed.FetchBytes(c.Url, content, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to fetch. url:%s", c.Url)
	}

	jst := time.FixedZone("Asia/Tokyo", 9 * 60 * 60)
	items := []*ChannelItem{}
	err = nil
	for _, ch := range feed.Channels {
		for _, item := range ch.Items {
			url := item.Links[0].Href
			publishedAt, err := time.Parse(
				"2006-01-02T15:04:05Z",
				item.PubDate,
			)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to parse time. pubDate:%v", item.PubDate)
			}
			publishedAt = publishedAt.UTC().In(jst)

			items = append(items, &ChannelItem{
				Title: sanitize.HTML(item.Title),
				Url: url,
				PublishedAt: &publishedAt,
			})
		}
	}
	return items, err
}
