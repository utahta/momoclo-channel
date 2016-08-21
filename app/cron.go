package app

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/utahta/momoclo-channel/crawler"
	"github.com/utahta/momoclo-channel/log"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/taskqueue"
	"google.golang.org/appengine/urlfetch"
)

type CronHandler struct {
	context context.Context
	Log     log.Logger
}

func (h *CronHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.context = appengine.NewContext(r)
	h.Log = log.NewGaeLogger(h.context)

	switch r.URL.Path {
	case "/cron/crawl":
		h.serveCrawl(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *CronHandler) serveCrawl(w http.ResponseWriter, r *http.Request) {
	h.Log.Infof("crawl start.")

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

	h.Log.Infof("crawl end.")
}

func (h *CronHandler) pushTweetQueue(params url.Values) {
	task := taskqueue.NewPOSTTask("/queue/tweet", params)
	_, err := taskqueue.Add(h.context, task, "queue-tweet")
	if err != nil {
		h.Log.Errorf("Failed to add taskqueue for tweet. error:%v", err)
	}
}

func (h *CronHandler) pushLineQueue(params url.Values) {
	task := taskqueue.NewPOSTTask("/queue/line", params)
	_, err := taskqueue.Add(h.context, task, "queue-line")
	if err != nil {
		h.Log.Errorf("Failed to add taskqueue for line. error:%v", err)
	}
}

func (h *CronHandler) crawlChannelClients() []*crawler.ChannelClient {
	return []*crawler.ChannelClient{
		crawler.NewTamaiBlogChannelClient(),
		crawler.NewMomotaBlogChannelClient(),
		crawler.NewAriyasuBlogChannelClient(),
		crawler.NewSasakiBlogChannelClient(),
		crawler.NewTakagiBlogChannelClient(),
		crawler.NewAeNewsChannelClient(),
		crawler.NewYoutubeChannelClient(),
	}
}
