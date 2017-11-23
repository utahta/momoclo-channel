package usecase

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/event"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/domain/service/eventtask"
)

type (
	// Tweet use case
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
func (use *Tweet) Do(params TweetParams) error {
	const errTag = "Tweet.Do failed"

	if len(params.TweetRequests) == 0 {
		return errors.Errorf("%v: invalid tweet requests", errTag)
	}

	res, err := use.tweeter.Tweet(params.TweetRequests[0])
	if err != nil {
		return errors.Wrap(err, errTag)
	}
	use.log.Infof("tweet: %v", params.TweetRequests[0])

	requests := params.TweetRequests[1:] // go to next tweet
	if len(requests) == 0 {
		use.log.Info("done!")
		return nil
	}
	requests[0].InReplyToStatusID = res.IDStr

	task := eventtask.NewTweets(requests)
	if err := use.taskQueue.Push(task); err != nil {
		return errors.Wrap(err, errTag)
	}
	return nil
}
