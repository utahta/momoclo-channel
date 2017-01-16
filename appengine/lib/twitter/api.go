package twitter

import (
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/utahta/go-atomicbool"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/appengine/model"
	"github.com/utahta/momoclo-crawler"
	"golang.org/x/net/context"
)

// Tweet text message
func TweetMessage(ctx context.Context, text string) error {
	if tweetDisabled() {
		return nil
	}

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
	errFlg := atomicbool.New(false)
	var wg sync.WaitGroup
	wg.Add(len(ch.Items))
	for _, item := range ch.Items {
		go func(ctx context.Context, item *crawler.ChannelItem) {
			defer wg.Done()

			if err := model.NewTweetItem(item).Put(ctx); err != nil {
				return
			}
			if _, err := model.PutLatestBlogPost(ctx, item.Url); err != nil {
				log.GaeLog(ctx).Error(err)
				// go on
			}

			if err := tweetChannelItem(ctx, ch.Title, item); err != nil {
				errFlg.Set(true)
				log.GaeLog(ctx).Error(err)
				return
			}
		}(ctx, item)
	}
	wg.Wait()

	if errFlg.Enabled() {
		return errors.New("Errors occured in twitter.TweetChannel.")
	}
	return nil
}

func tweetChannelItem(ctx context.Context, title string, item *crawler.ChannelItem) error {
	if tweetDisabled() {
		return nil
	}

	tw, err := newChannelClient(ctx)
	if err != nil {
		return err
	}

	if err := tw.TweetItem(title, item); err != nil {
		return err
	}
	return nil
}

// if true disable tweet
func tweetDisabled() bool {
	e := os.Getenv("DISABLE_TWEET")
	if e != "" {
		return true
	}
	return false
}
