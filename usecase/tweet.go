package usecase

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/event"
	"github.com/utahta/momoclo-channel/domain/model"
)

type (
	// Tweet tweet any message to twitter
	Tweet struct {
		log       core.Logger
		taskQueue event.TaskQueue
		tweeter   model.Tweeter
	}

	// TweetParams input parameters
	TweetParams struct {
		TweetRequests []model.TweetRequest
	}
)

// NewTweet returns Tweet use case
func NewTweet(log core.Logger, taskQueue event.TaskQueue, tweeter model.Tweeter) *Tweet {
	return &Tweet{
		log:       log,
		taskQueue: taskQueue,
		tweeter:   tweeter,
	}
}

// Do tweet
func (t *Tweet) Do(params TweetParams) error {
	const errTag = "Tweet.Do failed"

	if len(params.TweetRequests) == 0 {
		t.log.Errorf("%v: invalid tweet requests", errTag)
		return errors.New("invalid tweet requests")
	}

	res, err := t.tweeter.Tweet(params.TweetRequests[0])
	if err != nil {
		return errors.Wrap(err, errTag)
	}
	t.log.Infof("tweet: %v", params.TweetRequests[0])

	tweetRequests := params.TweetRequests[1:] // go to next tweet
	if len(tweetRequests) == 0 {
		t.log.Infof("done!")
		return nil
	}
	tweetRequests[0].InReplyToStatusID = res.IDStr

	task := event.Task{QueueName: "queue-tweet", Path: "/queue/tweet", Object: tweetRequests}
	if err := t.taskQueue.Push(task); err != nil {
		return errors.Wrap(err, errTag)
	}
	return nil
}
