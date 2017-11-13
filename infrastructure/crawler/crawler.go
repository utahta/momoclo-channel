package crawler

import (
	"context"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain/entity"
	"github.com/utahta/momoclo-channel/usecase"
	"github.com/utahta/momoclo-crawler"
	"google.golang.org/appengine/urlfetch"
)

type handler struct {
	ctx context.Context
}

// New returns usecase.Crawler that wraps momoclo-crawler
func New(ctx context.Context) usecase.Crawler {
	return &handler{ctx}
}

func (c *handler) Fetch(code string, maxItemNum int, latestURL string) ([]*usecase.CrawlItem, error) {
	const errTag = "client.Fetch failed"
	var (
		cli  *crawler.ChannelClient
		err  error
		opts = crawler.WithHTTPClient(urlfetch.Client(c.ctx))
	)

	switch code {
	case entity.LatestEntryCodeTamai:
		cli, err = crawler.NewTamaiBlogChannelClient(maxItemNum, latestURL, opts)
	case entity.LatestEntryCodeMomota:
		cli, err = crawler.NewMomotaBlogChannelClient(maxItemNum, latestURL, opts)
	case entity.LatestEntryCodeAriyasu:
		cli, err = crawler.NewAriyasuBlogChannelClient(maxItemNum, latestURL, opts)
	case entity.LatestEntryCodeSasaki:
		cli, err = crawler.NewSasakiBlogChannelClient(maxItemNum, latestURL, opts)
	case entity.LatestEntryCodeTakagi:
		cli, err = crawler.NewTakagiBlogChannelClient(maxItemNum, latestURL, opts)
	case entity.LatestEntryCodeHappyclo:
		cli, err = crawler.NewHappycloChannelClient(latestURL, opts)
	case entity.LatestEntryCodeAeNews:
		cli, err = crawler.NewAeNewsChannelClient(opts)
	case entity.LatestEntryCodeYoutube:
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

	var items = make([]*usecase.CrawlItem, len(channel.Items))
	for i, feed := range channel.Items {
		item := &usecase.CrawlItem{
			Title:       channel.Title,
			URL:         channel.Url,
			EntryTitle:  feed.Title,
			EntryURL:    feed.Url,
			PublishedAt: *feed.PublishedAt,
		}

		item.ImageURLs = make([]string, len(feed.Images))
		for i, image := range feed.Images {
			item.ImageURLs[i] = image.Url
		}

		item.VideoURLs = make([]string, len(feed.Videos))
		for i, video := range feed.Videos {
			item.VideoURLs[i] = video.Url
		}

		items[i] = item
	}
	return items, nil
}
