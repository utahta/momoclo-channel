package usecase

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/event"
	"github.com/utahta/momoclo-channel/event/eventtask"
	"github.com/utahta/momoclo-channel/log"
	"github.com/utahta/momoclo-channel/timeutil"
	"github.com/utahta/momoclo-channel/types"
)

type (
	// Remind use case
	Remind struct {
		log       log.Logger
		taskQueue event.TaskQueue
		repo      types.ReminderRepository
	}
)

// NewRemind returns Remind use case
func NewRemind(
	logger log.Logger,
	taskQueue event.TaskQueue,
	repo types.ReminderRepository) *Remind {
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
			eventtask.NewTweet(types.TweetRequest{Text: reminder.Text}),
			eventtask.NewLineBroadcast(types.LineNotifyMessage{Text: fmt.Sprintf("\n%s", reminder.Text)}),
		})
		r.log.Infof("remind: %#v", reminder)
	}
	return nil
}
