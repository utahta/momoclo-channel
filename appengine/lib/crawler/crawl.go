package crawler

import (
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/crawler"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

func Crawl(ctx context.Context) error {
	var workQueue = make(chan bool, 20)
	defer close(workQueue)

	glog := log.NewGaeLogger(ctx)
	clients := crawlChannelClients(ctx)

	var errCount int32 = 0
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
				atomic.AddInt32(&errCount, 1)
				glog.Error(err)
				return
			}

			q := NewQueueTask(glog)
			if err := q.PushTweet(ctx, ch); err != nil {
				atomic.AddInt32(&errCount, 1)
				glog.Error(err)
			}
			if err := q.PushLine(ctx, ch); err != nil {
				atomic.AddInt32(&errCount, 1)
				glog.Error(err)
			}
		}(ctx, cli)
	}
	wg.Wait()

	if errCount > 0 {
		return errors.Errorf("Errors occured in crawler.Crawl. errCount:%d", errCount)
	}
	return nil
}

func crawlChannelClients(ctx context.Context) []*crawler.ChannelClient {
	option := crawler.WithHTTPClient(urlfetch.Client(ctx))
	return []*crawler.ChannelClient{
		retrieveChannelClient(crawler.NewTamaiBlogChannelClient(1, option)),
		retrieveChannelClient(crawler.NewMomotaBlogChannelClient(1, option)),
		retrieveChannelClient(crawler.NewAriyasuBlogChannelClient(1, option)),
		retrieveChannelClient(crawler.NewSasakiBlogChannelClient(1, option)),
		retrieveChannelClient(crawler.NewTakagiBlogChannelClient(1, option)),
		retrieveChannelClient(crawler.NewAeNewsChannelClient(option)),
		retrieveChannelClient(crawler.NewYoutubeChannelClient(option)),
	}
}

func retrieveChannelClient(c *crawler.ChannelClient, err error) *crawler.ChannelClient {
	return c
}
