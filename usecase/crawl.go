package usecase

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain"
	"github.com/utahta/momoclo-channel/domain/entity"
	"github.com/utahta/momoclo-channel/domain/service/latestentry"
	"github.com/utahta/momoclo-channel/lib/config"
	"github.com/utahta/momoclo-channel/lib/log"
	"github.com/utahta/momoclo-channel/lib/timeutil"
	"github.com/utahta/momoclo-crawler"
	"golang.org/x/sync/errgroup"
	"google.golang.org/appengine/urlfetch"
)

type (
	// Crawl crawling use case
	Crawl struct {
		latestEntryRepo entity.LatestEntryRepository
	}
)

// NewCrawl returns Crawl use case
func NewCrawl(latestEntryRepo entity.LatestEntryRepository) *Crawl {
	return &Crawl{
		latestEntryRepo: latestEntryRepo,
	}
}

// Do crawls some sites
func (c *Crawl) Do(ctx context.Context) error {
	const errTag = "Crawl.Crawl failed"

	var workQueue = make(chan bool, 20)
	defer close(workQueue)

	clients := c.channelClients(ctx)
	eg := &errgroup.Group{}
	for _, cli := range clients {
		workQueue <- true
		cli := cli

		eg.Go(func() error {
			defer func() {
				<-workQueue
			}()

			ch, err := cli.Fetch()
			if err != nil {
				log.Error(ctx, err)
				return err
			}

			c.updateLatestEntry(ctx, ch)

			if err := c.PushTweet(ctx, ch); err != nil {
				log.Errorf(ctx, "%v: push tweet queue. err:%v", errTag, err)
			}
			if err := c.PushLine(ctx, ch); err != nil {
				log.Errorf(ctx, "%v: push line queue. err:%v", errTag, err)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return errors.Wrap(err, "Errors occurred in crawler.Crawl")
	}
	return nil
}

func (c *Crawl) channelClients(ctx context.Context) []*crawler.ChannelClient {
	option := crawler.WithHTTPClient(urlfetch.Client(ctx))
	clients := []*crawler.ChannelClient{
		c.retrieveChannelClient(crawler.NewTamaiBlogChannelClient(1, c.latestEntryRepo.GetTamaiURL(), option)),
		c.retrieveChannelClient(crawler.NewMomotaBlogChannelClient(1, c.latestEntryRepo.GetMomotaURL(), option)),
		c.retrieveChannelClient(crawler.NewAriyasuBlogChannelClient(1, c.latestEntryRepo.GetAriyasuURL(), option)),
		c.retrieveChannelClient(crawler.NewSasakiBlogChannelClient(1, c.latestEntryRepo.GetSasakiURL(), option)),
		c.retrieveChannelClient(crawler.NewTakagiBlogChannelClient(1, c.latestEntryRepo.GetTakagiURL(), option)),
		c.retrieveChannelClient(crawler.NewAeNewsChannelClient(option)),
		c.retrieveChannelClient(crawler.NewYoutubeChannelClient(option)),
	}

	now := timeutil.Now().In(config.JST)
	if (now.Weekday() == time.Sunday && now.Hour() == 16 && (now.Minute() >= 55 && now.Minute() <= 59)) ||
		(now.Hour() >= 8 && now.Hour() <= 23 && (now.Minute() == 0 || now.Minute() == 30)) {
		clients = append(clients, c.retrieveChannelClient(crawler.NewHappycloChannelClient(c.latestEntryRepo.GetHappycloURL(), option)))
	}

	return clients
}

func (c *Crawl) retrieveChannelClient(cli *crawler.ChannelClient, _ error) *crawler.ChannelClient {
	return cli
}

func (c *Crawl) updateLatestEntry(ctx context.Context, ch *crawler.Channel) {
	const errTag = "Crawl.updateLatestEntry failed"

	for _, item := range ch.Items {
		l, err := c.latestEntryRepo.FindByURL(item.Url)
		if err == domain.ErrNoSuchEntity {
			l, err = latestentry.Parse(item.Url)
			if err != nil {
				log.Warningf(ctx, "%v: parse url:%v err:%v", errTag, item.Url, err)
				continue
			}
		} else if err != nil {
			log.Errorf(ctx, "%v: FindByURL url:%v err:%v", errTag, item.Url, err)
			continue
		} else {
			if l.URL == item.Url {
				continue
			}
		}

		if err := c.latestEntryRepo.Save(l); err != nil {
			log.Warningf(ctx, "%v: put latest entry. err:%v", errTag, err)
			continue
		}
		break // first item equals latest item
	}
}
