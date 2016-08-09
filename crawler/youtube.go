package crawler

import (
	"time"
	"io"
	"io/ioutil"

	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/pkg/errors"
)

type YouTubeChannel struct {
	*Channel
}

func NewYoutubeChannel() *YouTubeChannel {
	return &YouTubeChannel{Channel: &Channel{Url: "https://www.youtube.com/feeds/videos.xml?channel_id=UC7pcEjI2U2vg6CqgbwIpjgg"}}
}

func FetchYoutube() ([]*ChannelItem, error) {
	return FetchParse(NewYoutubeChannel())
}

func (c *YouTubeChannel) Fetch() (io.ReadCloser, error) {
	return c.fetch(c.Url)
}

func (c *YouTubeChannel) Parse(r io.Reader) ([]*ChannelItem, error) {
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
				"2006-01-02T15:04:05-07:00",
				item.PubDate,
			)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to parse time. pubDate:%v", item.PubDate)
			}
			publishedAt = publishedAt.UTC().In(jst)

			items = append(items, &ChannelItem{
				Title: item.Title,
				Url: url,
				PublishedAt: &publishedAt,
			})
		}
	}
	return items, nil
}
