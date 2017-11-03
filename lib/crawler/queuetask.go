package crawler

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/utahta/momoclo-channel/lib/linenotify"
	"github.com/utahta/momoclo-channel/lib/twitter"
	"github.com/utahta/momoclo-channel/domain"
	"github.com/utahta/momoclo-crawler"
	"google.golang.org/appengine/taskqueue"
)

type QueueTask struct{}

// New queue task
func NewQueueTask() *QueueTask {
	return &QueueTask{}
}

// Push task to tweet queue
func (q *QueueTask) PushTweet(ctx context.Context, ch *crawler.Channel) error {
	for _, item := range ch.Items {
		if err := domain.NewTweetItem(item).Put(ctx); err != nil {
			continue
		}

		param := &twitter.ChannelParam{Title: ch.Title, Item: item}
		v, err := q.buildURLValues(param)
		if err != nil {
			return err
		}

		task := taskqueue.NewPOSTTask("/queue/tweet", v)
		if _, err := taskqueue.Add(ctx, task, "queue-tweet"); err != nil {
			return err
		}
	}
	return nil
}

// Push task to LINE queue
func (q *QueueTask) PushLine(ctx context.Context, ch *crawler.Channel) error {
	for _, item := range ch.Items {
		if err := domain.NewLineItem(item).Put(ctx); err != nil {
			continue
		}

		param := &linenotify.ChannelParam{Title: ch.Title, Item: item}
		v, err := q.buildURLValues(param)
		if err != nil {
			return err
		}

		task := taskqueue.NewPOSTTask("/queue/line", v)
		if _, err := taskqueue.Add(ctx, task, "queue-line"); err != nil {
			return err
		}
	}
	return nil
}

func (q *QueueTask) ParseURLValues(v url.Values, ch interface{}) error {
	if err := json.Unmarshal([]byte(v.Get("channel")), ch); err != nil {
		return err
	}
	return nil
}

func (q *QueueTask) buildURLValues(ch interface{}) (url.Values, error) {
	bin, err := json.Marshal(ch)
	if err != nil {
		return nil, err
	}
	return url.Values{"channel": {string(bin)}}, nil
}
