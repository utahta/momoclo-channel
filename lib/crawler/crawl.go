package crawler

import (
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/lib/config"
	"github.com/utahta/momoclo-channel/lib/log"
	"github.com/utahta/momoclo-channel/model/latestentry"
	"github.com/utahta/momoclo-crawler"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"google.golang.org/appengine/urlfetch"
)

var timeNow = time.Now

func Crawl(ctx context.Context) error {
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
			for _, item := range ch.Items {
				if _, err := latestentry.Repository.PutURL(ctx, item.Url); err != nil {
					log.Errorf(ctx, "Failed to put latest entry. err:%v", err)
					// go on
				}
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
	clients := []*crawler.ChannelClient{
		retrieveChannelClient(crawler.NewTamaiBlogChannelClient(1, latestentry.Repository.GetTamaiURL(ctx), option)),
		retrieveChannelClient(crawler.NewMomotaBlogChannelClient(1, latestentry.Repository.GetMomotaURL(ctx), option)),
		retrieveChannelClient(crawler.NewAriyasuBlogChannelClient(1, latestentry.Repository.GetAriyasuURL(ctx), option)),
		retrieveChannelClient(crawler.NewSasakiBlogChannelClient(1, latestentry.Repository.GetSasakiURL(ctx), option)),
		retrieveChannelClient(crawler.NewTakagiBlogChannelClient(1, latestentry.Repository.GetTakagiURL(ctx), option)),
		retrieveChannelClient(crawler.NewAeNewsChannelClient(option)),
		retrieveChannelClient(crawler.NewYoutubeChannelClient(option)),
	}
	now := timeNow().In(config.JST)

	// every week on Sunday, 16:55 <= now <= 17:59 || 20:00 <= now <= 20:59
	if now.Weekday() == time.Sunday && ((now.Hour() == 16 && now.Minute() >= 55) || now.Hour() == 17 || now.Hour() == 20) {
		clients = append(clients, retrieveChannelClient(crawler.NewHappycloChannelClient(latestentry.Repository.GetHappycloURL(ctx), option)))
	}

	return clients
}

func retrieveChannelClient(c *crawler.ChannelClient, _ error) *crawler.ChannelClient {
	return c
}
