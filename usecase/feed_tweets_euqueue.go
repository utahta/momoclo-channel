package usecase

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/event"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/domain/service/tweet"
)

type (
	// FeedTweetsEnqueue use case
	FeedTweetsEnqueue struct {
		log       core.Logger
		taskQueue event.TaskQueue
	}

	// FeedTweetsEnqueueParams input parameters
	FeedTweetsEnqueueParams struct {
		FeedItem model.FeedItem
	}
)

// NewFeedTweetsEnqueue returns EnqueueTweets use case
func NewFeedTweetsEnqueue(log core.Logger, taskqueue event.TaskQueue) *FeedTweetsEnqueue {
	return &FeedTweetsEnqueue{
		log:       log,
		taskQueue: taskqueue,
	}
}

// Do enqueue feed tweet item
func (t *FeedTweetsEnqueue) Do(params FeedTweetsEnqueueParams) error {
	const errTag = "EnqueueTweets.Do failed"

	feedTweets := tweet.ConvertFeedTweets(params.FeedItem)
	if len(feedTweets) == 0 {
		t.log.Errorf("%v: invalid enqueue tweets feedItem:%v", errTag, params.FeedItem)
		return errors.New("invalid enqueue tweets")
	}

	task := event.Task{QueueName: "queue-tweet", Path: "/queue/feed/tweet", Object: feedTweets}
	if err := t.taskQueue.Push(task); err != nil {
		return errors.Wrap(err, errTag)
	}
	return nil
}
