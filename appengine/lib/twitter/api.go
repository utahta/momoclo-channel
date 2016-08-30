package twitter

import (
	"fmt"
	"os"
	"time"

	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/crawler"
	"github.com/utahta/momoclo-channel/twitter"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

func newChannelClient(ctx context.Context) *twitter.ChannelClient {
	tw := twitter.NewChannelClient(
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
		os.Getenv("TWITTER_ACCESS_TOKEN"),
		os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
	)
	tw.Log = log.NewGaeLogger(ctx)
	tw.Api.HttpClient.Transport = &urlfetch.Transport{Context: ctx}
	return tw
}

func newMessageClient(ctx context.Context) *twitter.MessageClient {
	tw := twitter.NewMessageClient(
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
		os.Getenv("TWITTER_ACCESS_TOKEN"),
		os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
	)
	tw.Log = log.NewGaeLogger(ctx)
	tw.Api.HttpClient.Transport = &urlfetch.Transport{Context: ctx}
	return tw
}

func TweetChannelItem(ctx context.Context, title string, item *crawler.ChannelItem) {
	tw := newChannelClient(ctx)
	if err := tw.TweetItem(title, item); err != nil {
		tw.Log.Error(err)
		return
	}
}

func TweetUstream(ctx context.Context) {
	tw := newMessageClient(ctx)

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	t := time.Now().In(jst)
	if err := tw.Tweet(fmt.Sprintf("momocloTV が配信を開始しました\n%s", t.Format("from 2006/01/02 15:04:05"))); err != nil {
		tw.Log.Error(err)
		return
	}
}

func TweetText(ctx context.Context, text string) {
	tw := newMessageClient(ctx)
	if err := tw.Tweet(text); err != nil {
		tw.Log.Error(err)
		return
	}
}
