package twitter

import (
	"os"

	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/twitter"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

func newChannelClient(ctx context.Context) (*twitter.ChannelClient, error) {
	tw, err := twitter.NewChannelClient(
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
		os.Getenv("TWITTER_ACCESS_TOKEN"),
		os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
		twitter.WithHTTPTransport(&urlfetch.Transport{Context: ctx}),
		twitter.WithLogger(log.NewGaeLogger(ctx)),
	)
	if err != nil {
		return nil, err
	}
	return tw, nil
}

func newMessageClient(ctx context.Context) (*twitter.MessageClient, error) {
	tw, err := twitter.NewMessageClient(
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
		os.Getenv("TWITTER_ACCESS_TOKEN"),
		os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
		twitter.WithHTTPTransport(&urlfetch.Transport{Context: ctx}),
		twitter.WithLogger(log.NewGaeLogger(ctx)),
	)
	if err != nil {
		return nil, err
	}
	return tw, nil
}
