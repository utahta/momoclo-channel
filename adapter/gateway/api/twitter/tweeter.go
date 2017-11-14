package twitter

import (
	"context"
	"net/url"

	"github.com/pkg/errors"
	"github.com/utahta/go-twitter"
	"github.com/utahta/go-twitter/types"
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/lib/config"
	"google.golang.org/appengine/urlfetch"
)

type tweeter struct {
	*twitter.Client
	log core.Logger
}

// NewTweeter returns model.Tweeter that wraps go-twitter
func NewTweeter(ctx context.Context, log core.Logger) model.Tweeter {
	if config.C.Twitter.Disabled {
		return &nop{}
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
	return &tweeter{c, log}
}

// TweetMessage tweets text
func (c *tweeter) TweetMessage(text string) error {
	if _, err := c.Tweet(text, nil); err != nil {
		return errors.Wrap(err, "TweetMessage failed")
	}
	return nil
}

// TweetFeed tweets feed
func (c *tweeter) TweetFeed(item model.FeedTweet) (model.FeedTweetResult, error) {
	const errTag = "TweetFeed failed"

	var (
		tweets *types.Tweets
		err    error
	)
	if item.Text != "" {
		if len(item.ImageURLs) > 0 {
			tweets, err = c.TweetImageURLs(item.Text, item.ImageURLs, nil)
		} else if item.VideoURL != "" {
			tweets, err = c.TweetVideoURL(item.Text, item.VideoURL, "video/mp4", nil)
		} else {
			tweets, err = c.Tweet(item.Text, nil)
		}
	} else {
		v := url.Values{}
		if item.InReplyToStatusID != "" {
			v.Set("in_reply_to_status_id", item.InReplyToStatusID)
		}
		if len(item.ImageURLs) > 0 {
			tweets, err = c.TweetImageURLs("", item.ImageURLs, v)
		} else if item.VideoURL != "" {
			tweets, err = c.TweetVideoURL("", item.VideoURL, "video/mp4", v)
		}
	}

	if err != nil {
		return model.FeedTweetResult{}, errors.Wrap(err, errTag)
	}
	return model.FeedTweetResult{IDStr: tweets.IDStr}, nil
}
