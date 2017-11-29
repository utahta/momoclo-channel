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
		Requests []model.TweetRequest `validate:"min=1,dive"`
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

	if err := core.Validate(params); err != nil {
		return errors.Wrap(err, errTag)
	}

	res, err := use.tweeter.Tweet(params.Requests[0])
	if err != nil {
		return errors.Wrap(err, errTag)
	}
	use.log.Infof("tweet: %v", params.Requests[0])

	requests := params.Requests[1:] // go to next tweet
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
