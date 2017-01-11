package crawler

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/utahta/go-atomicbool"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-crawler"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

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

func retrieveChannelClient(c *crawler.ChannelClient, _ error) *crawler.ChannelClient {
	return c
}
