package usecase

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/crawler"
	"github.com/utahta/momoclo-channel/dao"
	"github.com/utahta/momoclo-channel/entity"
	"github.com/utahta/momoclo-channel/event"
	"github.com/utahta/momoclo-channel/event/eventtask"
	"github.com/utahta/momoclo-channel/log"
	"github.com/utahta/momoclo-channel/validator"
)

type (
	// EnqueueTweets use case
	EnqueueTweets struct {
		log        log.Logger
		taskQueue  event.TaskQueue
		transactor dao.Transactor
		repo       entity.TweetItemRepository
	}

	// EnqueueTweetsParams input parameters
	EnqueueTweetsParams struct {
		FeedItem crawler.FeedItem
	}
)

// NewEnqueueTweets returns EnqueueTweets use case
func NewEnqueueTweets(
	log log.Logger,
	taskQueue event.TaskQueue,
	transactor dao.Transactor,
	repo entity.TweetItemRepository) *EnqueueTweets {
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

	item := entity.NewTweetItem(
		params.FeedItem.UniqueURL(),
		params.FeedItem.EntryTitle,
		params.FeedItem.EntryURL,
		params.FeedItem.PublishedAt,
		params.FeedItem.ImageURLs,
		params.FeedItem.VideoURLs,
	)
	if use.repo.Exists(item.ID) {
		return nil // already enqueued
	}

	err := use.transactor.RunInTransaction(func(h dao.PersistenceHandler) error {
		repo := use.repo.Tx(h)
		if _, err := repo.Find(item.ID); err != dao.ErrNoSuchEntity {
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
