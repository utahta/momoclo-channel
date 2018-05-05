package usecase

import (
	"fmt"

	"context"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/entity"
	"github.com/utahta/momoclo-channel/event"
	"github.com/utahta/momoclo-channel/event/eventtask"
	"github.com/utahta/momoclo-channel/linenotify"
	"github.com/utahta/momoclo-channel/log"
	"github.com/utahta/momoclo-channel/timeutil"
	"github.com/utahta/momoclo-channel/twitter"
)

type (
	// Remind use case
	Remind struct {
		log       log.Logger
		taskQueue event.TaskQueue
		repo      entity.ReminderRepository
	}
)

// NewRemind returns Remind use case
func NewRemind(
	logger log.Logger,
	taskQueue event.TaskQueue,
	repo entity.ReminderRepository) *Remind {
	return &Remind{
		log:       logger,
		taskQueue: taskQueue,
		repo:      repo,
	}
}

// Do remind
func (r *Remind) Do(ctx context.Context) error {
	const errTag = "Remind.Do failed"

	reminders, err := r.repo.FindAll(ctx)
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
			if err := r.repo.Save(ctx, reminder); err != nil {
				r.log.Errorf(ctx, "%v: update reminder %v", errTag, reminder)
				// not return error
			}
		}

		r.taskQueue.PushMulti(ctx, []event.Task{
			eventtask.NewTweet(twitter.TweetRequest{Text: reminder.Text}),
			eventtask.NewLineBroadcast(linenotify.Message{Text: fmt.Sprintf("\n%s", reminder.Text)}),
		})
		r.log.Infof(ctx, "remind: %#v", reminder)
	}
	return nil
}
