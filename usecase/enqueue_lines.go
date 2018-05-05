package usecase

import (
	"context"

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
	// EnqueueLines use case
	EnqueueLines struct {
		log        log.Logger
		taskQueue  event.TaskQueue
		transactor dao.Transactor
		repo       entity.LineItemRepository
	}

	// EnqueueLinesParams input parameters
	EnqueueLinesParams struct {
		FeedItem crawler.FeedItem
	}
)

// NewEnqueueLines returns EnqueueLines use case
func NewEnqueueLines(
	log log.Logger,
	taskQueue event.TaskQueue,
	transactor dao.Transactor,
	repo entity.LineItemRepository) *EnqueueLines {
	return &EnqueueLines{
		log:        log,
		taskQueue:  taskQueue,
		transactor: transactor,
		repo:       repo,
	}
}

// Do converts feeds to line notify requests and enqueue it
func (use *EnqueueLines) Do(ctx context.Context, params EnqueueLinesParams) error {
	const errTag = "EnqueueLines.Do failed"

	if err := validator.Validate(params); err != nil {
		return errors.Wrap(err, errTag)
	}

	item := entity.NewLineItem(
		params.FeedItem.UniqueURL(),
		params.FeedItem.EntryTitle,
		params.FeedItem.EntryURL,
		params.FeedItem.PublishedAt,
		params.FeedItem.ImageURLs,
		params.FeedItem.VideoURLs,
	)
	if use.repo.Exists(ctx, item.ID) {
		return nil // already enqueued
	}

	err := use.transactor.RunInTransaction(ctx, func(ctx context.Context) error {
		if _, err := use.repo.Find(ctx, item.ID); err != dao.ErrNoSuchEntity {
			return err
		}
		return use.repo.Save(ctx, item)
	}, nil)
	if err != nil {
		use.log.Errorf("%v: enqueue lines feedItem:%v", errTag, params.FeedItem)
		return errors.Wrap(err, errTag)
	}

	messages := params.FeedItem.ToLineNotifyMessages()
	if len(messages) == 0 {
		use.log.Errorf("%v: invalid enqueue lines feedItem:%v", errTag, params.FeedItem)
		return errors.Errorf("%v: invalid enqueue line messages", errTag)
	}

	task := eventtask.NewLinesBroadcast(messages)
	if err := use.taskQueue.Push(ctx, task); err != nil {
		return errors.Wrap(err, errTag)
	}
	use.log.Infof("enqueue line messages:%#v", messages)

	return nil
}
