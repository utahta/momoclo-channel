package usecase

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/event"
	"github.com/utahta/momoclo-channel/event/eventtask"
	"github.com/utahta/momoclo-channel/log"
	"github.com/utahta/momoclo-channel/timeutil"
)

type (
	// CheckUstream use case
	CheckUstream struct {
		log       log.Logger
		taskQueue event.TaskQueue
		checker   model.UstreamStatusChecker
		repo      model.UstreamStatusRepository
	}
)

// NewCheckUstream returns CheckUstream use case
func NewCheckUstream(
	logger log.Logger,
	taskQueue event.TaskQueue,
	checker model.UstreamStatusChecker,
	repo model.UstreamStatusRepository) *CheckUstream {
	return &CheckUstream{
		log:       logger,
		taskQueue: taskQueue,
		checker:   checker,
		repo:      repo,
	}
}

// Do checks momocloTV live status
func (u *CheckUstream) Do() error {
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
			eventtask.NewTweet(
				model.TweetRequest{Text: fmt.Sprintf("momocloTV が配信を開始しました\n%s\nhttp://www.ustream.tv/channel/momoclotv", t.Format("from 2006/01/02 15:04:05"))},
			),
			eventtask.NewLineBroadcast(model.LineNotifyMessage{Text: "\nmomocloTV が配信を開始しました\nhttp://www.ustream.tv/channel/momoclotv"}),
		})
	}
	return nil
}
