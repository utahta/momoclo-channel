package crawler

import (
	"sync"

	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/crawler"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

func Crawl(ctx context.Context) error {
	var workQueue = make(chan bool, 20)
	defer close(workQueue)

	client := urlfetch.Client(ctx)
	glog := log.NewGaeLogger(ctx)

	var wg sync.WaitGroup
	for _, cli := range crawlChannelClients() {
		workQueue <- true
		wg.Add(1)
		go func(ctx context.Context, cli *crawler.ChannelClient) {
			defer func() {
				<-workQueue
				wg.Done()
			}()
			cli.Channel.Client = client

			ch, err := cli.Fetch()
			if err != nil {
				glog.Error(err)
				return
			}

			q := NewQueueTask(glog)
			if err := q.PushTweet(ctx, ch); err != nil {
				glog.Error(err)
			}
			if err := q.PushLine(ctx, ch); err != nil {
				glog.Error(err)
			}
		}(ctx, cli)
	}
	wg.Wait()

	return nil
}

func crawlChannelClients() []*crawler.ChannelClient {
	bopt := &crawler.BlogChannelParserOption{MaxItemNum: 1}
	return []*crawler.ChannelClient{
		crawler.NewTamaiBlogChannelClient(bopt),
		crawler.NewMomotaBlogChannelClient(bopt),
		crawler.NewAriyasuBlogChannelClient(bopt),
		crawler.NewSasakiBlogChannelClient(bopt),
		crawler.NewTakagiBlogChannelClient(bopt),
		crawler.NewAeNewsChannelClient(),
		crawler.NewYoutubeChannelClient(),
	}
}
