package usecase

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain"
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/event"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/domain/service/eventtask"
	"github.com/utahta/momoclo-channel/log"
)

type (
	// LineNotify use case
	LineNotify struct {
		log       log.Logger
		taskQueue event.TaskQueue
		notify    model.LineNotify
		repo      model.LineNotificationRepository
	}

	// LineNotifyParams input parameters
	LineNotifyParams struct {
		Request model.LineNotifyRequest
	}
)

// NewLineNotify returns LineNotify use case
func NewLineNotify(
	log log.Logger,
	taskQueue event.TaskQueue,
	notify model.LineNotify,
	repo model.LineNotificationRepository) *LineNotify {
	return &LineNotify{
		log:       log,
		taskQueue: taskQueue,
		notify:    notify,
		repo:      repo,
	}
}

// Do line notify
func (use *LineNotify) Do(params LineNotifyParams) error {
	const errTag = "LineNotify.Do failed"

	if err := core.Validate(params); err != nil {
		return errors.Wrap(err, errTag)
	}

	request := params.Request
	err := use.notify.Notify(request.AccessToken, request.Messages[0])
	if err != nil {
		if err == domain.ErrInvalidAccessToken {
			err = use.repo.Delete(request.ID)
			use.log.Infof("delete id:%v err:%v", request.ID, err)
		}
		return errors.Wrap(err, errTag)
	}
	use.log.Infof("line notify id:%v message:%v", request.ID, request.Messages[0])

	request.Messages = request.Messages[1:]
	if len(request.Messages) == 0 {
		use.log.Info("done!")
		return nil
	}

	if err := use.taskQueue.Push(eventtask.NewLine(request)); err != nil {
		return errors.Wrap(err, errTag)
	}
	return nil
}
