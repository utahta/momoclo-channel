package usecase

import (
	"context"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/event"
	"github.com/utahta/momoclo-channel/event/eventtask"
	"github.com/utahta/momoclo-channel/log"
	"github.com/utahta/momoclo-channel/twitter"
	"github.com/utahta/momoclo-channel/validator"
)

type (
	// Tweet use case
	Tweet struct {
		log       log.Logger
		taskQueue event.TaskQueue
		tweeter   twitter.Tweeter
	}

	// TweetParams input parameters
	TweetParams struct {
		Requests []twitter.TweetRequest `validate:"min=1,dive"`
	}
)

// NewTweet returns Tweet use case
func NewTweet(log log.Logger, taskQueue event.TaskQueue, tweeter twitter.Tweeter) *Tweet {
	return &Tweet{
		log:       log,
		taskQueue: taskQueue,
		tweeter:   tweeter,
	}
}

// Do tweet
func (use *Tweet) Do(ctx context.Context, params TweetParams) error {
	const errTag = "Tweet.Do failed"

	if err := validator.Validate(params); err != nil {
		return errors.Wrap(err, errTag)
	}

	res, err := use.tweeter.Tweet(ctx, params.Requests[0])
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
	if err := use.taskQueue.Push(ctx, task); err != nil {
		return errors.Wrap(err, errTag)
	}
	return nil
}
