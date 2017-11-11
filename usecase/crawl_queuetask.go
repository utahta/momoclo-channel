package usecase

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/utahta/momoclo-channel/domain"
	"github.com/utahta/momoclo-channel/lib/linenotify"
	"github.com/utahta/momoclo-channel/lib/twitter"
	"github.com/utahta/momoclo-crawler"
	"google.golang.org/appengine/taskqueue"
)

//FIXME

// PushTweet queues tweet task
func (c *Crawl) PushTweet(ctx context.Context, ch *crawler.Channel) error {
	for _, item := range ch.Items {
		if err := domain.NewTweetItem(item).Put(ctx); err != nil {
			continue
		}

		param := &twitter.ChannelParam{Title: ch.Title, Item: item}
		v, err := c.buildURLValues(param)
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

// PushLine queues line notification task
func (c *Crawl) PushLine(ctx context.Context, ch *crawler.Channel) error {
	for _, item := range ch.Items {
		if err := domain.NewLineItem(item).Put(ctx); err != nil {
			continue
		}

		param := &linenotify.ChannelParam{Title: ch.Title, Item: item}
		v, err := c.buildURLValues(param)
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

// ParseURLValues parses channel params
func (c *Crawl) ParseURLValues(v url.Values, ch interface{}) error {
	if err := json.Unmarshal([]byte(v.Get("channel")), ch); err != nil {
		return err
	}
	return nil
}

func (c *Crawl) buildURLValues(ch interface{}) (url.Values, error) {
	bin, err := json.Marshal(ch)
	if err != nil {
		return nil, err
	}
	return url.Values{"channel": {string(bin)}}, nil
}
