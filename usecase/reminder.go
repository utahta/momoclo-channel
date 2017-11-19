package usecase

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/lib/timeutil"
	"golang.org/x/sync/errgroup"
)

type (
	// Reminder use case
	Reminder struct {
		log     core.Logger
		repo    model.ReminderRepository
		tweeter model.Tweeter
	}
)

// NewReminder returns Reminder use case
func NewReminder(
	logger core.Logger,
	repo model.ReminderRepository,
	tweeter model.Tweeter) *Reminder {
	return &Reminder{
		log:     logger,
		repo:    repo,
		tweeter: tweeter,
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

		var eg errgroup.Group
		eg.Go(func() error {
			if _, err := r.tweeter.Tweet(model.TweetRequest{Text: reminder.Text}); err != nil {
				r.log.Error("%v: tweet", errTag)
				return err
			}
			return nil
		})

		eg.Go(func() error {
			//FIXME
			//if err := linenotify.NotifyMessage(ctx, row.Text); err != nil {
			//	logger.Error(ctx, err)
			//	return err
			//}
			return nil
		})

		if reminder.IsOneTime() {
			reminder.Disable()
			if err := r.repo.Save(reminder); err != nil {
				r.log.Errorf("%v: update reminder %v", errTag, reminder)
				// not return error
			}
		}

		if err := eg.Wait(); err != nil {
			r.log.Errorf("%v: remind text:%#v", reminder)
			return errors.Wrap(err, errTag)
		}
	}
	return nil
}
