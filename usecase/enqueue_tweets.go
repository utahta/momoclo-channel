package usecase

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/event"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/domain/service/tweet"
)

type (
	// EnqueueTweets use case
	EnqueueTweets struct {
		log       core.Logger
		taskQueue event.TaskQueue
	}

	// EnqueueTweetsParams input parameters
	EnqueueTweetsParams struct {
		FeedItem model.FeedItem
	}
)

// NewEnqueueTweets returns EnqueueTweets use case
func NewEnqueueTweets(log core.Logger, taskqueue event.TaskQueue) *EnqueueTweets {
	return &EnqueueTweets{
		log:       log,
		taskQueue: taskqueue,
	}
}

// Do converts feeds to tweet requests and enqueue it
func (t *EnqueueTweets) Do(params EnqueueTweetsParams) error {
	const errTag = "EnqueueTweets.Do failed"

	tweetRequests := tweet.ConvertToTweetRequests(params.FeedItem)
	if len(tweetRequests) == 0 {
		t.log.Errorf("%v: invalid enqueue tweets feedItem:%v", errTag, params.FeedItem)
		return errors.New("invalid enqueue tweets")
	}

	task := event.Task{QueueName: "queue-tweet", Path: "/queue/tweet", Object: tweetRequests}
	if err := t.taskQueue.Push(task); err != nil {
		return errors.Wrap(err, errTag)
	}
	return nil
}
