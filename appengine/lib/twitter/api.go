package twitter

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/appengine/model"
	"github.com/utahta/momoclo-channel/crawler"
	"golang.org/x/net/context"
)

// Tweet text message
func TweetMessage(ctx context.Context, text string) error {
	tw, err := newMessageClient(ctx)
	if err != nil {
		return err
	}

	if err := tw.Tweet(text); err != nil {
		return err
	}
	return nil
}

// Tweet channel
func TweetChannel(ctx context.Context, ch *crawler.Channel) error {
	errs := make([]error, len(ch.Items))
	var wg sync.WaitGroup
	wg.Add(len(ch.Items))
	for i, item := range ch.Items {
		go func(ctx context.Context, item *crawler.ChannelItem, i int) {
			defer wg.Done()

			if err := model.NewTweetItem(item).Put(ctx); err != nil {
				return
			}
			if err := tweetChannelItem(ctx, ch.Title, item); err != nil {
				errs[i] = err
				return
			}
		}(ctx, item, i)
	}
	wg.Wait()

	any := false
	for _, err := range errs {
		if err != nil {
			any = true
			log.GaeLog(ctx).Error(err)
		}
	}
	if any {
		return errors.New("Errors occured in twitter.TweetChannel.")
	}
	return nil
}

func tweetChannelItem(ctx context.Context, title string, item *crawler.ChannelItem) error {
	tw, err := newChannelClient(ctx)
	if err != nil {
		return err
	}

	if err := tw.TweetItem(title, item); err != nil {
		return err
	}
	return nil
}
