package usecase

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/entity"
	"github.com/utahta/momoclo-channel/event"
	"github.com/utahta/momoclo-channel/event/eventtask"
	"github.com/utahta/momoclo-channel/linenotify"
	"github.com/utahta/momoclo-channel/log"
	"github.com/utahta/momoclo-channel/validator"
)

type (
	// LineNotify use case
	LineNotify struct {
		log       log.Logger
		taskQueue event.TaskQueue
		notify    linenotify.Client
		repo      entity.LineNotificationRepository
	}

	// LineNotifyParams input parameters
	LineNotifyParams struct {
		Request linenotify.Request
	}
)

// NewLineNotify returns LineNotify use case
func NewLineNotify(
	log log.Logger,
	taskQueue event.TaskQueue,
	notify linenotify.Client,
	repo entity.LineNotificationRepository) *LineNotify {
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

	if err := validator.Validate(params); err != nil {
		return errors.Wrap(err, errTag)
	}

	request := params.Request
	err := use.notify.Notify(request.AccessToken, request.Messages[0])
	if err != nil {
		if err == linenotify.ErrInvalidAccessToken {
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
