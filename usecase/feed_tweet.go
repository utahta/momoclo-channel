package usecase

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/event"
	"github.com/utahta/momoclo-channel/domain/model"
)

type (
	// FeedTweet use case
	FeedTweet struct {
		log       core.Logger
		taskQueue event.TaskQueue
		tweeter   model.Tweeter
	}

	// FeedTweetParams input parameters
	FeedTweetParams struct {
		FeedTweets []model.FeedTweet
	}
)

// NewFeedTweet returns FeedTweet use case
func NewFeedTweet(log core.Logger, taskqueue event.TaskQueue, tweeter model.Tweeter) *FeedTweet {
	return &FeedTweet{
		log:       log,
		taskQueue: taskqueue,
		tweeter:   tweeter,
	}
}

// Do enqueue tweets feed
func (t *FeedTweet) Do(params FeedTweetParams) error {
	const errTag = "FeedTweet.Do failed"

	if len(params.FeedTweets) == 0 {
		t.log.Errorf("%v: invalid feed tweets", errTag)
		return errors.New("invalid feed tweets")
	}

	res, err := t.tweeter.TweetFeed(params.FeedTweets[0])
	if err != nil {
		return errors.Wrap(err, errTag)
	}
	t.log.Infof("feed tweeted: %v", params.FeedTweets[0])

	feedTweets := params.FeedTweets[1:]
	if len(feedTweets) == 0 {
		t.log.Infof("feed tweet done!")
		return nil
	}
	feedTweets[0].InReplyToStatusID = res.IDStr

	task := event.Task{QueueName: "queue-tweet", Path: "/queue/feed/tweet", Object: feedTweets}
	if err := t.taskQueue.Push(task); err != nil {
		return errors.Wrap(err, errTag)
	}
	return nil
}
