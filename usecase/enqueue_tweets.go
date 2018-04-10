package usecase

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/event"
	"github.com/utahta/momoclo-channel/event/eventtask"
	"github.com/utahta/momoclo-channel/log"
	"github.com/utahta/momoclo-channel/types"
	"github.com/utahta/momoclo-channel/validator"
)

type (
	// EnqueueTweets use case
	EnqueueTweets struct {
		log        log.Logger
		taskQueue  event.TaskQueue
		transactor types.Transactor
		repo       types.TweetItemRepository
	}

	// EnqueueTweetsParams input parameters
	EnqueueTweetsParams struct {
		FeedItem types.FeedItem
	}
)

// NewEnqueueTweets returns EnqueueTweets use case
func NewEnqueueTweets(
	log log.Logger,
	taskQueue event.TaskQueue,
	transactor types.Transactor,
	repo types.TweetItemRepository) *EnqueueTweets {
	return &EnqueueTweets{
		log:        log,
		taskQueue:  taskQueue,
		transactor: transactor,
		repo:       repo,
	}
}

// Do converts feeds to tweet requests and enqueue it
func (use *EnqueueTweets) Do(params EnqueueTweetsParams) error {
	const errTag = "EnqueueTweets.Do failed"

	if err := validator.Validate(params); err != nil {
		return errors.Wrap(err, errTag)
	}

	item := types.NewTweetItem(params.FeedItem)
	if use.repo.Exists(item.ID) {
		return nil // already enqueued
	}

	err := use.transactor.RunInTransaction(func(h types.PersistenceHandler) error {
		repo := use.repo.Tx(h)
		if _, err := repo.Find(item.ID); err != types.ErrNoSuchEntity {
			return err
		}
		return repo.Save(item)
	}, nil)
	if err != nil {
		use.log.Errorf("%v: enqueue tweets feedItem:%v", errTag, params.FeedItem)
		return errors.Wrap(err, errTag)
	}

	requests := params.FeedItem.ToTweetRequests()
	if len(requests) == 0 {
		use.log.Errorf("%v: invalid enqueue tweets feedItem:%v", errTag, params.FeedItem)
		return errors.Errorf("%v: invalid enqueue tweets", errTag)
	}

	task := eventtask.NewTweets(requests)
	if err := use.taskQueue.Push(task); err != nil {
		return errors.Wrap(err, errTag)
	}
	use.log.Infof("enqueue tweet requests:%#v", requests)

	return nil
}
