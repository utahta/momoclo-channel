package crawler

import (
	"context"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/types"
	"github.com/utahta/momoclo-crawler"
	"google.golang.org/appengine/urlfetch"
)

type handler struct {
	ctx context.Context
}

// New returns model.FeedFetcher that wraps momoclo-crawler
func New(ctx context.Context) types.FeedFetcher {
	return &handler{ctx}
}

func (c *handler) Fetch(code types.FeedCode, maxItemNum int, latestURL string) ([]types.FeedItem, error) {
	const errTag = "client.Fetch failed"
	var (
		cli  *crawler.ChannelClient
		err  error
		opts = crawler.WithHTTPClient(urlfetch.Client(c.ctx))
	)

	switch code {
	case types.FeedCodeTamai:
		cli, err = crawler.NewTamaiBlogChannelClient(maxItemNum, latestURL, opts)
	case types.FeedCodeMomota:
		cli, err = crawler.NewMomotaBlogChannelClient(maxItemNum, latestURL, opts)
	case types.FeedCodeAriyasu:
		cli, err = crawler.NewAriyasuBlogChannelClient(maxItemNum, latestURL, opts)
	case types.FeedCodeSasaki:
		cli, err = crawler.NewSasakiBlogChannelClient(maxItemNum, latestURL, opts)
	case types.FeedCodeTakagi:
		cli, err = crawler.NewTakagiBlogChannelClient(maxItemNum, latestURL, opts)
	case types.FeedCodeHappyclo:
		cli, err = crawler.NewHappycloChannelClient(latestURL, opts)
	case types.FeedCodeAeNews:
		cli, err = crawler.NewAeNewsChannelClient(opts)
	case types.FeedCodeYoutube:
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

	var items = make([]types.FeedItem, len(channel.Items))
	for i, feed := range channel.Items {
		item := types.FeedItem{
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
