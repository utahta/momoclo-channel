package usecase

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain"
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/event"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/lib/timeutil"
)

type (
	// CheckUstreamStatus use case
	CheckUstreamStatus struct {
		log       core.Logger
		taskQueue event.TaskQueue
		checker   model.UstreamStatusChecker
		repo      model.UstreamStatusRepository
	}
)

// NewCheckUstreamStatus returns CheckUstreamStatus use case
func NewCheckUstreamStatus(
	logger core.Logger,
	taskQueue event.TaskQueue,
	checker model.UstreamStatusChecker,
	repo model.UstreamStatusRepository) *CheckUstreamStatus {
	return &CheckUstreamStatus{
		log:       logger,
		taskQueue: taskQueue,
		checker:   checker,
		repo:      repo,
	}
}

// Do checks ustream status
func (u *CheckUstreamStatus) Do() error {
	const errTag = "CheckUstream.Do failed"

	isLive, err := u.checker.IsLive()
	if err != nil {
		return errors.Wrap(err, errTag)
	}

	status, err := u.repo.Find(model.UstreamStatusID)
	if err != nil && err != domain.ErrNoSuchEntity {
		return errors.Wrap(err, errTag)
	}
	if status.IsLive == isLive {
		return nil // nothing to do
	}

	status.IsLive = isLive
	if err := u.repo.Save(status); err != nil {
		return errors.Wrap(err, errTag)
	}

	if isLive {
		t := timeutil.Now()
		u.taskQueue.PushMulti([]event.Task{
			{QueueName: "queue-tweet", Path: "/queue/tweet", Object: []model.TweetRequest{
				{Text: fmt.Sprintf("momocloTV が配信を開始しました\n%s\nhttp://www.ustream.tv/channel/momoclotv", t.Format("from 2006/01/02 15:04:05"))},
			}},
			//FIXME add line event
		})
	}
	return nil
}
