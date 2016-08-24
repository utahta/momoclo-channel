package app

import (
	"encoding/json"
	"net/url"
	"sync"
	"time"

	"github.com/utahta/momoclo-channel/crawler"
	"github.com/utahta/momoclo-channel/log"
	"golang.org/x/net/context"
	"google.golang.org/appengine/taskqueue"
	"google.golang.org/appengine/urlfetch"
)

type Crawler struct {
	context context.Context
	log     log.Logger
}

func newCrawler(ctx context.Context) *Crawler {
	return &Crawler{context: ctx, log: log.NewGaeLogger(ctx)}
}

func (c *Crawler) Crawl() *Error {
	var workQueue = make(chan bool, 20)
	defer close(workQueue)

	ctx, cancel := context.WithTimeout(c.context, 50*time.Second)
	defer cancel()
	client := urlfetch.Client(ctx)

	var wg sync.WaitGroup
	for _, cli := range c.crawlChannelClients() {
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
				c.log.Errorf("Failed to fetch. error:%v", err)
				return
			}

			bin, err := json.Marshal(ch)
			if err != nil {
				c.log.Errorf("Failed to encode to json. error:%v", err)
				return
			}
			params := url.Values{"channel": {string(bin)}}

			c.pushTweetQueue(ctx, params)
			c.pushLineQueue(ctx, params)
		}(ctx, cli)
	}
	wg.Wait()

	return nil
}

func (c *Crawler) pushTweetQueue(ctx context.Context, params url.Values) {
	task := taskqueue.NewPOSTTask("/queue/tweet", params)
	_, err := taskqueue.Add(ctx, task, "queue-tweet")
	if err != nil {
		c.log.Errorf("Failed to add taskqueue for tweet. error:%v", err)
	}
}

func (c *Crawler) pushLineQueue(ctx context.Context, params url.Values) {
	task := taskqueue.NewPOSTTask("/queue/line", params)
	_, err := taskqueue.Add(ctx, task, "queue-line")
	if err != nil {
		c.log.Errorf("Failed to add taskqueue for line. error:%v", err)
	}
}

func (c *Crawler) crawlChannelClients() []*crawler.ChannelClient {
	bopt := &crawler.BlogChannelParserOption{MaxItemNum: 2}
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
