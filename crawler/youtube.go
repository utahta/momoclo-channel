package crawler

import (
	"io"
	"io/ioutil"
	"time"

	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/pkg/errors"
)

type youtubeChannelParser struct {
	channel *Channel
}

func NewYoutubeChannelClient(options ...ChannelClientOption) (*ChannelClient, error) {
	c := newChannel("https://www.youtube.com/feeds/videos.xml?channel_id=UC7pcEjI2U2vg6CqgbwIpjgg", "ニュータイプ放送局")
	return newChannelClient(c, &youtubeChannelParser{channel: c}, options...)
}

func FetchYoutube() (*Channel, error) {
	cc, err := NewYoutubeChannelClient()
	if err != nil {
		return nil, err
	}
	return cc.Fetch()
}

func (p *youtubeChannelParser) Parse(r io.Reader) ([]*ChannelItem, error) {
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

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	items := []*ChannelItem{}
	for _, ch := range feed.Channels {
		for _, item := range ch.Items {
			url := item.Links[0].Href
			publishedAt, err := time.Parse(
				"2006-01-02T15:04:05-07:00",
				item.PubDate,
			)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to parse time. pubDate:%v", item.PubDate)
			}
			publishedAt = publishedAt.UTC().In(jst)

			items = append(items, &ChannelItem{
				Title:       item.Title,
				Url:         url,
				PublishedAt: &publishedAt,
			})
		}
	}
	return items, nil
}
