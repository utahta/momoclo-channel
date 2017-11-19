package usecase

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain"
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/event"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/domain/service/convert"
)

type (
	// EnqueueTweets use case
	EnqueueTweets struct {
		log        core.Logger
		taskQueue  event.TaskQueue
		transactor model.Transactor
		repo       model.TweetItemRepository
	}

	// EnqueueTweetsParams input parameters
	EnqueueTweetsParams struct {
		FeedItem model.FeedItem
	}
)

// NewEnqueueTweets returns EnqueueTweets use case
func NewEnqueueTweets(
	log core.Logger,
	taskQueue event.TaskQueue,
	transactor model.Transactor,
	repo model.TweetItemRepository) *EnqueueTweets {
	return &EnqueueTweets{
		log:        log,
		taskQueue:  taskQueue,
		transactor: transactor,
		repo:       repo,
	}
}

// Do converts feeds to tweet requests and enqueue it
func (t *EnqueueTweets) Do(params EnqueueTweetsParams) error {
	const errTag = "EnqueueTweets.Do failed"

	tweetItem := convert.FeedItemToTweetItem(params.FeedItem)
	if t.repo.Exists(tweetItem.ID) {
		return nil // already enqueued
	}

	err := t.transactor.RunInTransaction(func(h model.PersistenceHandler) error {
		done := t.transactor.With(h, t.repo)
		defer done()

		if _, err := t.repo.Find(tweetItem.ID); err != domain.ErrNoSuchEntity {
			return err
		}
		return t.repo.Save(tweetItem)
	}, nil)
	if err != nil {
		t.log.Errorf("%v: enqueue tweets feedItem:%v", errTag, params.FeedItem)
		return errors.Wrap(err, errTag)
	}

	tweetRequests := convert.FeedItemToTweetRequests(params.FeedItem)
	if len(tweetRequests) == 0 {
		t.log.Errorf("%v: invalid enqueue tweets feedItem:%v", errTag, params.FeedItem)
		return errors.New("invalid enqueue tweets")
	}

	task := event.Task{QueueName: "queue-tweet", Path: "/queue/tweet", Object: tweetRequests}
	if err := t.taskQueue.Push(task); err != nil {
		return errors.Wrap(err, errTag)
	}
	t.log.Infof("enqueue tweet requests:%#v", tweetRequests)

	return nil
}
