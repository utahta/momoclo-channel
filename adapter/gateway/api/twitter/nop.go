package twitter

import (
	"github.com/utahta/momoclo-channel/domain/model"
)

type nop struct{}

func (c *nop) TweetMessage(text string) error {
	return nil
}

func (c *nop) TweetFeed(item model.FeedTweet) (model.FeedTweetResult, error) {
	return model.FeedTweetResult{}, nil
}
