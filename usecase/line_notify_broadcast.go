package usecase

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/config"
	"github.com/utahta/momoclo-channel/event"
	"github.com/utahta/momoclo-channel/event/eventtask"
	"github.com/utahta/momoclo-channel/log"
	"github.com/utahta/momoclo-channel/types"
	"github.com/utahta/momoclo-channel/validator"
)

type (
	// LineNotifyBroadcast use case
	LineNotifyBroadcast struct {
		log       log.Logger
		taskQueue event.TaskQueue
		repo      types.LineNotificationRepository
	}

	// LineNotifyBroadcastParams input parameters
	LineNotifyBroadcastParams struct {
		Messages []types.LineNotifyMessage `validate:"min=1,dive"`
	}
)

// NewLineNotifyBroadcast returns LineNotifyBroadcast use case
func NewLineNotifyBroadcast(
	log log.Logger,
	taskQueue event.TaskQueue,
	repo types.LineNotificationRepository) *LineNotifyBroadcast {
	return &LineNotifyBroadcast{
		log:       log,
		taskQueue: taskQueue,
		repo:      repo,
	}
}

// Do notify broadcast
func (use *LineNotifyBroadcast) Do(params LineNotifyBroadcastParams) error {
	const errTag = "LineNotifyBroadcast.Do failed"

	if err := validator.Validate(params); err != nil {
		return errors.Wrap(err, errTag)
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
		tasks = append(tasks, eventtask.NewLine(types.LineNotifyRequest{
			ID:          n.ID,
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
