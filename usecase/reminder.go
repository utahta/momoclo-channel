package usecase

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/event"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/lib/timeutil"
)

type (
	// Reminder use case
	Reminder struct {
		log       core.Logger
		taskQueue event.TaskQueue
		repo      model.ReminderRepository
	}
)

// NewReminder returns Reminder use case
func NewReminder(
	logger core.Logger,
	taskQueue event.TaskQueue,
	repo model.ReminderRepository) *Reminder {
	return &Reminder{
		log:       logger,
		taskQueue: taskQueue,
		repo:      repo,
	}
}

// Do remind
func (r *Reminder) Do() error {
	const errTag = "Reminder.Do failed"

	reminders, err := r.repo.FindAll()
	if err != nil {
		return errors.Wrap(err, errTag)
	}

	currentTime := timeutil.Now()
	for _, reminder := range reminders {
		if ok, err := reminder.Valid(currentTime); !ok {
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
			{QueueName: "queue-tweet", Path: "/queue/tweet", Object: []model.TweetRequest{{Text: reminder.Text}}},
			//FIXME add line event
		})
		r.log.Infof("remind: %#v", reminder)
	}
	return nil
}
