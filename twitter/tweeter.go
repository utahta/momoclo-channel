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
		Tweet(context.Context, TweetRequest) (TweetResponse, error)
	}

	tweeter struct {
	}
)

// NewTweeter returns model.Tweeter that wraps go-twitter
func NewTweeter() Tweeter {
	if config.C().Twitter.Disabled {
		return NewNopTweeter()
	}

	twitter.SetConsumerCredentials(
		config.C().Twitter.ConsumerKey,
		config.C().Twitter.ConsumerSecret,
	)
	return &tweeter{}
}

// Tweet tweets given request
func (t *tweeter) Tweet(ctx context.Context, req TweetRequest) (TweetResponse, error) {
	const errTag = "tweeter.Tweet failed"

	c, err := twitter.New(
		config.C().Twitter.AccessToken,
		config.C().Twitter.AccessTokenSecret,
		twitter.WithHTTPClient(urlfetch.Client(ctx)),
	)
	if err != nil {
		return TweetResponse{}, errors.Wrap(err, errTag)
	}

	var tweets *types.Tweets
	if req.Text != "" {
		if len(req.ImageURLs) > 0 {
			tweets, err = c.TweetImageURLs(req.Text, req.ImageURLs, nil)
		} else if req.VideoURL != "" {
			tweets, err = c.TweetVideoURL(req.Text, req.VideoURL, "video/mp4", nil)
		} else {
			tweets, err = c.Tweet(req.Text, nil)
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
