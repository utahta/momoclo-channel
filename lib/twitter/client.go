package twitter

import (
	"context"

	"github.com/utahta/go-twitter"
	"github.com/utahta/momoclo-channel/lib/config"
	"google.golang.org/appengine/urlfetch"
)

func newClient(ctx context.Context) (*twitter.Client, error) {
	twitter.SetConsumerCredentials(
		config.C.Twitter.ConsumerKey,
		config.C.Twitter.ConsumerSecret,
	)
	c, err := twitter.New(
		config.C.Twitter.AccessToken,
		config.C.Twitter.AccessTokenSecret,
		twitter.WithHTTPClient(urlfetch.Client(ctx)),
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}
