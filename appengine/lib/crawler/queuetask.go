package crawler

import (
	"encoding/json"
	"net/url"

	"github.com/utahta/momoclo-channel/appengine/lib/linenotify"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/appengine/lib/twitter"
	"github.com/utahta/momoclo-channel/appengine/model"
	"github.com/utahta/momoclo-channel/crawler"
	"golang.org/x/net/context"
	"google.golang.org/appengine/taskqueue"
	"sync"
)

type QueueTask struct {
	Log log.Logger
}

// New queue task
func NewQueueTask(log log.Logger) *QueueTask {
	return &QueueTask{
		Log: log,
	}
}

// Push task to tweet queue
func (q *QueueTask) PushTweet(ctx context.Context, ch *crawler.Channel) error {
	v, err := q.buildURLValues(ch)
	if err != nil {
		return err
	}

	task := taskqueue.NewPOSTTask("/queue/tweet", v)
	if _, err := taskqueue.Add(ctx, task, "queue-tweet"); err != nil {
		return err
	}
	return nil
}

// Push task to LINE queue
func (q *QueueTask) PushLine(ctx context.Context, ch *crawler.Channel) error {
	v, err := q.buildURLValues(ch)
	if err != nil {
		return err
	}

	task := taskqueue.NewPOSTTask("/queue/line", v)
	if _, err := taskqueue.Add(ctx, task, "queue-line"); err != nil {
		return err
	}
	return nil
}

// Run tweet task
func (q *QueueTask) RunTweet(ctx context.Context, v url.Values) error {
	ch, err := q.parseURLValues(v)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(ch.Items))
	for _, item := range ch.Items {
		go func(ctx context.Context, item *crawler.ChannelItem) {
			defer wg.Done()

			if err := model.NewTweetItem(item).Put(ctx); err != nil {
				return
			}
			twitter.TweetChannelItem(ctx, ch.Title, item)
		}(ctx, item)
	}
	wg.Wait()
	return nil
}

// Run LINE task
func (q *QueueTask) RunLine(ctx context.Context, v url.Values) error {
	ch, err := q.parseURLValues(v)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(ch.Items))
	for _, item := range ch.Items {
		go func(ctx context.Context, item *crawler.ChannelItem) {
			defer wg.Done()

			if err := model.NewLineItem(item).Put(ctx); err != nil {
				return
			}
			linenotify.NotifyChannelItem(ctx, ch.Title, item)
		}(ctx, item)
	}
	wg.Wait()
	return nil
}

func (q *QueueTask) parseURLValues(v url.Values) (*crawler.Channel, error) {
	var ch crawler.Channel
	if err := json.Unmarshal([]byte(v.Get("channel")), &ch); err != nil {
		return nil, err
	}
	return &ch, nil
}

func (q *QueueTask) buildURLValues(ch *crawler.Channel) (url.Values, error) {
	bin, err := json.Marshal(ch)
	if err != nil {
		return nil, err
	}
	return url.Values{"channel": {string(bin)}}, nil
}
