package app

import (
	"net/http"
	"net/url"
	"encoding/json"
	"sync"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/appengine/taskqueue"
	"google.golang.org/appengine/log"
	"github.com/utahta/momoclo-channel/crawler"
	"golang.org/x/net/context"
)

type CronHandler struct {
	context context.Context
}

func (h *CronHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.context = appengine.NewContext(r)

	switch r.URL.Path {
	case "/cron/crawl":
		h.serveCrawl(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *CronHandler) serveCrawl(w http.ResponseWriter, r *http.Request) {
	log.Infof(h.context, "crawl start.")

	var workQueue = make(chan bool, 5)
	defer close(workQueue)

	var wg sync.WaitGroup
	httpClient := urlfetch.Client(h.context)
	for _, c := range h.crawlChannelClients() {
		c.Channel.HttpClient = httpClient

		workQueue <- true
		wg.Add(1)
		go func(ctx context.Context, c *crawler.ChannelClient) {
			defer func(){
				<-workQueue
				wg.Done()
			}()

			items, err := c.Fetch()
			if err != nil {
				log.Errorf(ctx, "Failed to fetch. error:%v", err)
				return
			}

			bin, err := json.Marshal(items)
			if err != nil {
				log.Errorf(ctx, "Failed to encode to json. error:%v", err)
				return
			}
			params := url.Values{ "items": {string(bin)} }

			h.addTweetPushQueue(params)
			h.addLinePullQueue(params)
		}(h.context, c)
	}
	wg.Wait()

	log.Infof(h.context, "crawl end.")
}

func (h *CronHandler) addTweetPushQueue(params url.Values) {
	task := taskqueue.NewPOSTTask("/queue/tweet", params)
	_, err := taskqueue.Add(h.context, task, "queue-tweet")
	if err != nil {
		log.Errorf(h.context, "Failed to add taskqueue for tweet. error:%v", err)
	}
}

func (h *CronHandler) addLinePullQueue(params url.Values) {
	task := &taskqueue.Task{
		Payload: []byte(params.Encode()),
		Method:  "PULL",
	}
	_, err := taskqueue.Add(h.context, task, "queue-line")
	if err != nil {
		log.Errorf(h.context, "Failed to add taskqueue for line. error:%v", err)
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
