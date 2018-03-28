package usecase

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain/service/feeditem"
	"github.com/utahta/momoclo-channel/event"
	"github.com/utahta/momoclo-channel/event/eventtask"
	"github.com/utahta/momoclo-channel/log"
	"github.com/utahta/momoclo-channel/types"
	"github.com/utahta/momoclo-channel/validator"
)

type (
	// EnqueueLines use case
	EnqueueLines struct {
		log        log.Logger
		taskQueue  event.TaskQueue
		transactor types.Transactor
		repo       types.LineItemRepository
	}

	// EnqueueLinesParams input parameters
	EnqueueLinesParams struct {
		FeedItem types.FeedItem
	}
)

// NewEnqueueLines returns EnqueueLines use case
func NewEnqueueLines(
	log log.Logger,
	taskQueue event.TaskQueue,
	transactor types.Transactor,
	repo types.LineItemRepository) *EnqueueLines {
	return &EnqueueLines{
		log:        log,
		taskQueue:  taskQueue,
		transactor: transactor,
		repo:       repo,
	}
}

// Do converts feeds to line notify requests and enqueue it
func (use *EnqueueLines) Do(params EnqueueLinesParams) error {
	const errTag = "EnqueueLines.Do failed"

	if err := validator.Validate(params); err != nil {
		return errors.Wrap(err, errTag)
	}

	item := types.NewLineItem(params.FeedItem)
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
		use.log.Errorf("%v: enqueue lines feedItem:%v", errTag, params.FeedItem)
		return errors.Wrap(err, errTag)
	}

	messages := feeditem.ToLineNotifyMessages(params.FeedItem)
	if len(messages) == 0 {
		use.log.Errorf("%v: invalid enqueue lines feedItem:%v", errTag, params.FeedItem)
		return errors.Errorf("%v: invalid enqueue line messages", errTag)
	}

	task := eventtask.NewLinesBroadcast(messages)
	if err := use.taskQueue.Push(task); err != nil {
		return errors.Wrap(err, errTag)
	}
	use.log.Infof("enqueue line messages:%#v", messages)

	return nil
}
