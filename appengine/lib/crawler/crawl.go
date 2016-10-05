package crawler

import (
	"encoding/json"
	"net/url"
	"sync"

	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/crawler"
	"golang.org/x/net/context"
	"google.golang.org/appengine/taskqueue"
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
				glog.Errorf("Failed to fetch. error:%+v", err)
				return
			}

			bin, err := json.Marshal(ch)
			if err != nil {
				glog.Errorf("Failed to encode to json. error:%+v", err)
				return
			}
			params := url.Values{"channel": {string(bin)}}

			pushTweetQueue(ctx, params)
			pushLineQueue(ctx, params)
		}(ctx, cli)
	}
	wg.Wait()

	return nil
}

func pushTweetQueue(ctx context.Context, params url.Values) {
	task := taskqueue.NewPOSTTask("/queue/tweet", params)
	_, err := taskqueue.Add(ctx, task, "queue-tweet")
	if err != nil {
		log.GaeLog(ctx).Errorf("Failed to add taskqueue for tweet. error:%+v", err)
	}
}

func pushLineQueue(ctx context.Context, params url.Values) {
	task := taskqueue.NewPOSTTask("/queue/line", params)
	_, err := taskqueue.Add(ctx, task, "queue-line")
	if err != nil {
		log.GaeLog(ctx).Errorf("Failed to add taskqueue for line. error:%+v", err)
	}
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
