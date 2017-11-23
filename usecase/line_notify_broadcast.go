package usecase

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/event"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/domain/service/eventtask"
	"github.com/utahta/momoclo-channel/lib/config"
)

type (
	// LineNotifyBroadcast use case
	LineNotifyBroadcast struct {
		log       core.Logger
		taskQueue event.TaskQueue
		repo      model.LineNotificationRepository
	}

	// LineNotifyBroadcastParams input parameters
	LineNotifyBroadcastParams struct {
		Messages []model.LineNotifyMessage
	}
)

// NewLineNotifyBroadcast returns LineNotifyBroadcast use case
func NewLineNotifyBroadcast(
	log core.Logger,
	taskQueue event.TaskQueue,
	repo model.LineNotificationRepository) *LineNotifyBroadcast {
	return &LineNotifyBroadcast{
		log:       log,
		taskQueue: taskQueue,
		repo:      repo,
	}
}

// Do notify broadcast
func (use *LineNotifyBroadcast) Do(params LineNotifyBroadcastParams) error {
	const errTag = "LineNotifyBroadcast.Do failed"

	if len(params.Messages) == 0 {
		return errors.Errorf("%v: invalid line notify messages", errTag)
	}

	//TODO use iterator
	ns, err := use.repo.FindAll()
	if err != nil {
		return errors.Wrap(err, errTag)
	}

	tasks := make([]event.Task, 0, len(ns))
	for _, n := range ns {
		accessToken, err := n.Token(config.C.LineNotify.TokenKey)
		if err != nil {
			use.log.Errorf("%v: get access token err:%v", errTag, err)
			continue
		}
		tasks = append(tasks, eventtask.NewLine(model.LineNotifyRequest{
			AccessToken: accessToken,
			Messages:    params.Messages,
		}))
	}

	if err := use.taskQueue.PushMulti(tasks); err != nil {
		return errors.Wrap(err, errTag)
	}
	use.log.Infof("broadcast line tasks len:%v", len(tasks))

	return nil
}
