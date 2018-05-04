package twitter

import (
	"context"
	"net/url"

	"github.com/pkg/errors"
	"github.com/utahta/go-twitter"
	"github.com/utahta/go-twitter/types"
	"github.com/utahta/momoclo-channel/config"
	"google.golang.org/appengine/urlfetch"
)

type (
	// TweetRequest represents request that tweet message, img urls and video url data
	TweetRequest struct {
		InReplyToStatusID string
		Text              string
		ImageURLs         []string `validate:"dive,omitempty,url"`
		VideoURL          string   `validate:"omitempty,url"`
	}

	// TweetResponse represents response tweet data
	TweetResponse struct {
		IDStr string
	}

	// Tweeter interface
	Tweeter interface {
		Tweet(TweetRequest) (TweetResponse, error)
	}

	tweeter struct {
		*twitter.Client
	}
)

// NewTweeter returns model.Tweeter that wraps go-twitter
func NewTweeter(ctx context.Context) Tweeter {
	if config.C().Twitter.Disabled {
		return NewNopTweeter()
	}

	twitter.SetConsumerCredentials(
		config.C().Twitter.ConsumerKey,
		config.C().Twitter.ConsumerSecret,
	)
	c, _ := twitter.New(
		config.C().Twitter.AccessToken,
		config.C().Twitter.AccessTokenSecret,
		twitter.WithHTTPClient(urlfetch.Client(ctx)),
	)
	return &tweeter{c}
}

// Tweet tweets given request
func (c *tweeter) Tweet(req TweetRequest) (TweetResponse, error) {
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
		return TweetResponse{}, errors.Wrap(err, errTag)
	}
	return TweetResponse{IDStr: tweets.IDStr}, nil
}
