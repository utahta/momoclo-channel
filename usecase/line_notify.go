package usecase

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/event"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/domain/service/eventtask"
)

type (
	// LineNotify use case
	LineNotify struct {
		log       core.Logger
		taskQueue event.TaskQueue
		notify    model.LineNotify
	}

	// LineNotifyParams input parameters
	LineNotifyParams struct {
		Request model.LineNotifyRequest
	}
)

// NewLineNotify returns LineNotify use case
func NewLineNotify(
	log core.Logger,
	taskQueue event.TaskQueue,
	notify model.LineNotify) *LineNotify {
	return &LineNotify{
		log:       log,
		taskQueue: taskQueue,
		notify:    notify,
	}
}

// Do line notify
func (use *LineNotify) Do(params LineNotifyParams) error {
	const errTag = "LineNotify.Do failed"

	request := params.Request
	if len(request.Messages) == 0 {
		return errors.Errorf("%v: invalid line notify request", errTag)
	}

	err := use.notify.Notify(request.AccessToken, request.Messages[0])
	if err != nil {
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
