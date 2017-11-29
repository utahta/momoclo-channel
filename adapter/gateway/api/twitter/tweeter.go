package twitter

import (
	"context"
	"net/url"

	"github.com/pkg/errors"
	"github.com/utahta/go-twitter"
	"github.com/utahta/go-twitter/types"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/lib/config"
	"google.golang.org/appengine/urlfetch"
)

type tweeter struct {
	*twitter.Client
}

// NewTweeter returns model.Tweeter that wraps go-twitter
func NewTweeter(ctx context.Context) model.Tweeter {
	if config.C.Twitter.Disabled {
		return NewNopTweeter()
	}

	twitter.SetConsumerCredentials(
		config.C.Twitter.ConsumerKey,
		config.C.Twitter.ConsumerSecret,
	)
	c, _ := twitter.New(
		config.C.Twitter.AccessToken,
		config.C.Twitter.AccessTokenSecret,
		twitter.WithHTTPClient(urlfetch.Client(ctx)),
	)
	return &tweeter{c}
}

// Tweet tweets given request
func (c *tweeter) Tweet(req model.TweetRequest) (model.TweetResponse, error) {
	const errTag = "tweeter.Tweet failed"

	var (
		tweets *types.Tweets
		err    error
	)
	if req.Text != "" {
		if len(req.ImageURLs) > 0 {
			tweets, err = c.TweetImageURLs(req.Text, req.ImageURLs, nil)
		} else if req.VideoURL != "" {
			tweets, err = c.TweetVideoURL(req.Text, req.VideoURL, "video/mp4", nil)
		} else {
			tweets, err = c.Client.Tweet(req.Text, nil)
		}
	} else {
		v := url.Values{}
		if req.InReplyToStatusID != "" {
			v.Set("in_reply_to_status_id", req.InReplyToStatusID)
		}
		if len(req.ImageURLs) > 0 {
			tweets, err = c.TweetImageURLs("", req.ImageURLs, v)
		} else if req.VideoURL != "" {
			tweets, err = c.TweetVideoURL("", req.VideoURL, "video/mp4", v)
		}
	}

	if err != nil {
		return model.TweetResponse{}, errors.Wrap(err, errTag)
	}
	return model.TweetResponse{IDStr: tweets.IDStr}, nil
}
