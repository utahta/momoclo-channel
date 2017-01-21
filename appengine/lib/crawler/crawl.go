package crawler

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/go-atomicbool"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/appengine/model"
	"github.com/utahta/momoclo-crawler"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

var timeNow = time.Now

func Crawl(ctx context.Context) error {
	var workQueue = make(chan bool, 20)
	defer close(workQueue)

	glog := log.NewGaeLogger(ctx)
	clients := crawlChannelClients(ctx)

	errFlg := atomicbool.New(false)
	var wg sync.WaitGroup
	wg.Add(len(clients))
	for _, cli := range clients {
		workQueue <- true
		go func(ctx context.Context, cli *crawler.ChannelClient) {
			defer func() {
				<-workQueue
				wg.Done()
			}()

			ch, err := cli.Fetch()
			if err != nil {
				errFlg.Set(true)
				glog.Error(err)
				return
			}

			// update latest entry
			for _, item := range ch.Items {
				if _, err := model.PutLatestEntry(ctx, item.Url); err != nil {
					glog.Error(err)
					// go on
				}
			}

			q := NewQueueTask(glog)
			if err := q.PushTweet(ctx, ch); err != nil {
				errFlg.Set(true)
				glog.Error(err)
			}
			if err := q.PushLine(ctx, ch); err != nil {
				errFlg.Set(true)
				glog.Error(err)
			}
		}(ctx, cli)
	}
	wg.Wait()

	if errFlg.Enabled() {
		return errors.New("Errors occured in crawler.Crawl.")
	}
	return nil
}

func crawlChannelClients(ctx context.Context) []*crawler.ChannelClient {
	option := crawler.WithHTTPClient(urlfetch.Client(ctx))
	clients := []*crawler.ChannelClient{
		retrieveChannelClient(crawler.NewTamaiBlogChannelClient(1, model.GetTamaiLatestEntryURL(ctx), option)),
		retrieveChannelClient(crawler.NewMomotaBlogChannelClient(1, model.GetMomotaLatestEntryURL(ctx), option)),
		retrieveChannelClient(crawler.NewAriyasuBlogChannelClient(1, model.GetAriyasuLatestEntryURL(ctx), option)),
		retrieveChannelClient(crawler.NewSasakiBlogChannelClient(1, model.GetSasakiLatestEntryURL(ctx), option)),
		retrieveChannelClient(crawler.NewTakagiBlogChannelClient(1, model.GetTakagiLatestEntryURL(ctx), option)),
		retrieveChannelClient(crawler.NewAeNewsChannelClient(option)),
		retrieveChannelClient(crawler.NewYoutubeChannelClient(option)),
	}
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := timeNow().In(jst)

	// every week on Sunday, 16:55 <= now <= 19:00
	if now.Weekday() == time.Sunday && ((now.Hour() == 16 && now.Minute() >= 55) || now.Hour() >= 17) && now.Hour() <= 19 {
		clients = append(clients, retrieveChannelClient(crawler.NewHappycloChannelClient(model.GetHappycloLatestEntryURL(ctx), option)))
	}

	return clients
}

func retrieveChannelClient(c *crawler.ChannelClient, _ error) *crawler.ChannelClient {
	return c
}
