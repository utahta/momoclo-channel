package crawler

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/adapter/persistence"
	"github.com/utahta/momoclo-channel/domain"
	"github.com/utahta/momoclo-channel/domain/service/latestentry"
	"github.com/utahta/momoclo-channel/infrastructure/datastore"
	"github.com/utahta/momoclo-channel/lib/config"
	"github.com/utahta/momoclo-channel/lib/log"
	"github.com/utahta/momoclo-channel/lib/timeutil"
	"github.com/utahta/momoclo-crawler"
	"golang.org/x/sync/errgroup"
	"google.golang.org/appengine/urlfetch"
)

func Crawl(ctx context.Context) error {
	const errTag = "Crawl failed"
	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	var workQueue = make(chan bool, 20)
	defer close(workQueue)

	clients := crawlChannelClients(ctx)
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

			// update latest entry
			repo := persistence.NewLatestEntryRepository(datastore.New(ctx))
			for _, item := range ch.Items {
				l, err := repo.FindByURL(item.Url)
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

				if err := repo.Save(l); err != nil {
					log.Warningf(ctx, "%v: put latest entry. err:%v", errTag, err)
					continue
				}
				break // first item is the latest item
			}

			q := NewQueueTask()
			if err := q.PushTweet(ctx, ch); err != nil {
				log.Errorf(ctx, "Failed to push tweet queue. err:%v", err)
			}
			if err := q.PushLine(ctx, ch); err != nil {
				log.Errorf(ctx, "Failed to push line queue. err:%v", err)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return errors.Wrap(err, "Errors occurred in crawler.Crawl")
	}
	return nil
}

func crawlChannelClients(ctx context.Context) []*crawler.ChannelClient {
	option := crawler.WithHTTPClient(urlfetch.Client(ctx))
	repo := persistence.NewLatestEntryRepository(datastore.New(ctx))
	clients := []*crawler.ChannelClient{
		retrieveChannelClient(crawler.NewTamaiBlogChannelClient(1, repo.GetTamaiURL(), option)),
		retrieveChannelClient(crawler.NewMomotaBlogChannelClient(1, repo.GetMomotaURL(), option)),
		retrieveChannelClient(crawler.NewAriyasuBlogChannelClient(1, repo.GetAriyasuURL(), option)),
		retrieveChannelClient(crawler.NewSasakiBlogChannelClient(1, repo.GetSasakiURL(), option)),
		retrieveChannelClient(crawler.NewTakagiBlogChannelClient(1, repo.GetTakagiURL(), option)),
		retrieveChannelClient(crawler.NewAeNewsChannelClient(option)),
		retrieveChannelClient(crawler.NewYoutubeChannelClient(option)),
	}

	now := timeutil.Now().In(config.JST)
	if (now.Weekday() == time.Sunday && now.Hour() == 16 && (now.Minute() >= 55 && now.Minute() <= 59)) ||
		(now.Hour() >= 8 && now.Hour() <= 23 && (now.Minute() == 0 || now.Minute() == 30)) {
		clients = append(clients, retrieveChannelClient(crawler.NewHappycloChannelClient(repo.GetHappycloURL(), option)))
	}

	return clients
}

func retrieveChannelClient(c *crawler.ChannelClient, _ error) *crawler.ChannelClient {
	return c
}
