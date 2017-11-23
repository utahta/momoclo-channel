package usecase

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/event"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/domain/service/eventtask"
	"github.com/utahta/momoclo-channel/lib/timeutil"
)

type (
	// Remind use case
	Remind struct {
		log       core.Logger
		taskQueue event.TaskQueue
		repo      model.ReminderRepository
	}
)

// NewRemind returns Remind use case
func NewRemind(
	logger core.Logger,
	taskQueue event.TaskQueue,
	repo model.ReminderRepository) *Remind {
	return &Remind{
		log:       logger,
		taskQueue: taskQueue,
		repo:      repo,
	}
}

// Do remind
func (r *Remind) Do() error {
	const errTag = "Remind.Do failed"

	reminders, err := r.repo.FindAll()
	if err != nil {
		return errors.Wrap(err, errTag)
	}

	now := timeutil.Now()
	for _, reminder := range reminders {
		if ok, err := reminder.Valid(now); !ok {
			if err != nil {
				return err
			}
			continue
		}

		if reminder.IsOneTime() {
			reminder.Disable()
			if err := r.repo.Save(reminder); err != nil {
				r.log.Errorf("%v: update reminder %v", errTag, reminder)
				// not return error
			}
		}

		r.taskQueue.PushMulti([]event.Task{
			eventtask.NewTweet(model.TweetRequest{Text: reminder.Text}),
			eventtask.NewLineBroadcast(model.LineNotifyMessage{Text: reminder.Text}),
		})
		r.log.Infof("remind: %#v", reminder)
	}
	return nil
}
