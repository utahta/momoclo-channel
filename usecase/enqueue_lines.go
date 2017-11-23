package usecase

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain"
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/event"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/domain/service/eventtask"
	"github.com/utahta/momoclo-channel/domain/service/feeditem"
)

type (
	// EnqueueLines use case
	EnqueueLines struct {
		log        core.Logger
		taskQueue  event.TaskQueue
		transactor model.Transactor
		repo       model.LineItemRepository
	}

	// EnqueueLinesParams input parameters
	EnqueueLinesParams struct {
		FeedItem model.FeedItem
	}
)

// NewEnqueueLines returns EnqueueLines use case
func NewEnqueueLines(
	log core.Logger,
	taskQueue event.TaskQueue,
	transactor model.Transactor,
	repo model.LineItemRepository) *EnqueueLines {
	return &EnqueueLines{
		log:        log,
		taskQueue:  taskQueue,
		transactor: transactor,
		repo:       repo,
	}
}

// Do converts feeds to line notify requests and enqueue it
func (t *EnqueueLines) Do(params EnqueueLinesParams) error {
	const errTag = "EnqueueLines.Do failed"

	item := model.NewLineItem(params.FeedItem)
	if t.repo.Exists(item.ID) {
		return nil // already enqueued
	}

	err := t.transactor.RunInTransaction(func(h model.PersistenceHandler) error {
		done := t.transactor.With(h, t.repo)
		defer done()

		if _, err := t.repo.Find(item.ID); err != domain.ErrNoSuchEntity {
			return err
		}
		return t.repo.Save(item)
	}, nil)
	if err != nil {
		t.log.Errorf("%v: enqueue lines feedItem:%v", errTag, params.FeedItem)
		return errors.Wrap(err, errTag)
	}

	messages := feeditem.ToLineNotifyMessages(params.FeedItem)
	if len(messages) == 0 {
		t.log.Errorf("%v: invalid enqueue lines feedItem:%v", errTag, params.FeedItem)
		return errors.New("invalid enqueue line messages")
	}

	task := eventtask.NewLinesBroadcast(messages)
	if err := t.taskQueue.Push(task); err != nil {
		return errors.Wrap(err, errTag)
	}
	t.log.Infof("enqueue line messages:%#v", messages)

	return nil
}
