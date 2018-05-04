package crawler

import (
	"context"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-crawler"
	"google.golang.org/appengine/urlfetch"
)

type (
	// FeedFetcher interface
	FeedFetcher interface {
		Fetch(code FeedCode, maxItemNum int, latestURL string) ([]FeedItem, error)
	}

	client struct {
		ctx context.Context
	}
)

// New returns FeedFetcher that wraps momoclo-crawler
func New(ctx context.Context) FeedFetcher {
	return &client{ctx}
}

func (c *client) Fetch(code FeedCode, maxItemNum int, latestURL string) ([]FeedItem, error) {
	const errTag = "crawler Fetch failed"
	var (
		cli  *crawler.ChannelClient
		err  error
		opts = crawler.WithHTTPClient(urlfetch.Client(c.ctx))
	)

	switch code {
	case FeedCodeTamai:
		cli, err = crawler.NewTamaiBlogChannelClient(maxItemNum, latestURL, opts)
	case FeedCodeMomota:
		cli, err = crawler.NewMomotaBlogChannelClient(maxItemNum, latestURL, opts)
	case FeedCodeSasaki:
		cli, err = crawler.NewSasakiBlogChannelClient(maxItemNum, latestURL, opts)
	case FeedCodeTakagi:
		cli, err = crawler.NewTakagiBlogChannelClient(maxItemNum, latestURL, opts)
	case FeedCodeHappyclo:
		cli, err = crawler.NewHappycloChannelClient(latestURL, opts)
	case FeedCodeAeNews:
		cli, err = crawler.NewAeNewsChannelClient(opts)
	case FeedCodeYoutube:
		cli, err = crawler.NewYoutubeChannelClient(opts)
	default:
		err = errors.Errorf("code:%s did not support", code)
	}

	if err != nil {
		return nil, errors.Wrap(err, errTag)
	}

	channel, err := cli.Fetch()
	if err != nil {
		return nil, errors.Wrap(err, errTag)
	}

	var items = make([]FeedItem, len(channel.Items))
	for i, feed := range channel.Items {
		item := FeedItem{
			Title:       channel.Title,
			URL:         channel.URL,
			EntryTitle:  feed.Title,
			EntryURL:    feed.URL,
			PublishedAt: feed.PublishedAt,
		}

		item.ImageURLs = make([]string, len(feed.Images))
		for i, image := range feed.Images {
			item.ImageURLs[i] = image.URL
		}

		item.VideoURLs = make([]string, len(feed.Videos))
		for i, video := range feed.Videos {
			item.VideoURLs[i] = video.URL
		}

		items[i] = item
	}
	return items, nil
}
