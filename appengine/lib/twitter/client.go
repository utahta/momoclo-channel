package twitter

import (
	"os"

	"github.com/utahta/go-twitter"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

func newClient(ctx context.Context) (*twitter.Client, error) {
	twitter.SetConsumerCredentials(
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
	)
	c, err := twitter.New(
		os.Getenv("TWITTER_ACCESS_TOKEN"),
		os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
		twitter.WithHTTPClient(urlfetch.Client(ctx)),
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}
