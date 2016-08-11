package momoclo_channel

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

func crawlHandler(w http.ResponseWriter, r *http.Request) *appError {
	ctx := appengine.NewContext(r)
	log.Infof(ctx, "crawl start.")

	var workQueue = make(chan bool, 5)
	defer close(workQueue)

	var wg sync.WaitGroup
	httpClient := urlfetch.Client(ctx)
	for _, c := range crawlChannelClients() {
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

			addTweetPushQueue(ctx, params)
			addLinePullQueue(ctx, params)
		}(ctx, c)
	}
	wg.Wait()

	log.Infof(ctx, "crawl end.")
	return nil
}

func addTweetPushQueue(ctx context.Context, params url.Values) {
	task := taskqueue.NewPOSTTask("/queue/tweet", params)
	_, err := taskqueue.Add(ctx, task, "queue-tweet")
	if err != nil {
		log.Errorf(ctx, "Failed to add taskqueue for tweet. error:%v", err)
	}
}

func addLinePullQueue(ctx context.Context, params url.Values) {
	task := &taskqueue.Task{
		Payload: []byte(params.Encode()),
		Method:  "PULL",
	}
	_, err := taskqueue.Add(ctx, task, "queue-line")
	if err != nil {
		log.Errorf(ctx, "Failed to add taskqueue for line. error:%v", err)
	}
}

func crawlChannelClients() []*crawler.ChannelClient {
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
