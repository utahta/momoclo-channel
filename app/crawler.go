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
	Log     log.Logger
}

func (h *Crawler) Crawl(ctx context.Context) error {
	h.context = ctx
	h.Log = log.NewGaeLogger(h.context)

	var workQueue = make(chan bool, 20)
	defer close(workQueue)

	var wg sync.WaitGroup
	for _, c := range h.crawlChannelClients() {
		workQueue <- true
		wg.Add(1)
		go func(ctx context.Context, c *crawler.ChannelClient) {
			defer func() {
				<-workQueue
				wg.Done()
			}()

			ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()
			c.Channel.Client = urlfetch.Client(ctx)

			ch, err := c.Fetch()
			if err != nil {
				h.Log.Errorf("Failed to fetch. error:%v", err)
				return
			}

			bin, err := json.Marshal(ch)
			if err != nil {
				h.Log.Errorf("Failed to encode to json. error:%v", err)
				return
			}
			params := url.Values{"channel": {string(bin)}}

			h.pushTweetQueue(params)
			h.pushLineQueue(params)
		}(h.context, c)
	}
	wg.Wait()

	return nil
}

func (h *Crawler) pushTweetQueue(params url.Values) {
	task := taskqueue.NewPOSTTask("/queue/tweet", params)
	_, err := taskqueue.Add(h.context, task, "queue-tweet")
	if err != nil {
		h.Log.Errorf("Failed to add taskqueue for tweet. error:%v", err)
	}
}

func (h *Crawler) pushLineQueue(params url.Values) {
	task := taskqueue.NewPOSTTask("/queue/line", params)
	_, err := taskqueue.Add(h.context, task, "queue-line")
	if err != nil {
		h.Log.Errorf("Failed to add taskqueue for line. error:%v", err)
	}
}

func (h *Crawler) crawlChannelClients() []*crawler.ChannelClient {
	return []*crawler.ChannelClient{
		crawler.NewTamaiBlogChannelClient(nil),
		crawler.NewMomotaBlogChannelClient(nil),
		crawler.NewAriyasuBlogChannelClient(nil),
		crawler.NewSasakiBlogChannelClient(nil),
		crawler.NewTakagiBlogChannelClient(nil),
		crawler.NewAeNewsChannelClient(),
		crawler.NewYoutubeChannelClient(),
	}
}
