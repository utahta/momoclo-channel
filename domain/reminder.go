package domain

import (
	"context"
	"time"

	"github.com/mjibson/goon"
	"github.com/pkg/errors"
)

const (
	ReminderOnce ReminderType = iota
	ReminderWeekly
)

type ReminderType int

type Reminder struct {
	Id      int64 `datastore:"-" goon:"id"`
	Text    string
	Type    ReminderType
	Enabled bool

	// Once
	RemindAt time.Time `datastore:",noindex"`

	// Weekly
	Weekday time.Weekday `datastore:",noindex"`
	Hour    int          `datastore:",noindex"`
	Minute  int          `datastore:",noindex"`
}

func NewReminderOnce(text string, remindAt time.Time) *Reminder {
	return &Reminder{
		Text:     text,
		Type:     ReminderOnce,
		Enabled:  true,
		RemindAt: remindAt,
	}
}

func NewReminderWeekly(text string, weekday time.Weekday, hour int, minute int) *Reminder {
	return &Reminder{
		Text:    text,
		Type:    ReminderWeekly,
		Enabled: true,
		Weekday: weekday,
		Hour:    hour,
		Minute:  minute,
	}
}

func (r *Reminder) Put(ctx context.Context) error {
	g := goon.FromContext(ctx)
	_, err := g.Put(r)
	return err
}

func (r *Reminder) Disable(ctx context.Context) error {
	g := goon.FromContext(ctx)
	r.Enabled = false
	_, err := g.Put(r)
	return err
}

func (r *Reminder) IsOnce() bool {
	return r.Type == ReminderOnce
}

func (r *Reminder) Valid(now time.Time) (bool, error) {
	switch r.Type {
	case ReminderOnce:
		r.RemindAt = r.RemindAt.In(now.Location())
		if r.RemindAt.Year() == now.Year() && r.RemindAt.Month() == now.Month() && r.RemindAt.Day() == now.Day() &&
			r.RemindAt.Hour() == now.Hour() && r.RemindAt.Minute() == now.Minute() {
			return true, nil
		}
	case ReminderWeekly:
		if r.Weekday == now.Weekday() && r.Hour == now.Hour() && r.Minute == now.Minute() {
			return true, nil
		}
	default:
		return false, errors.Errorf("Invalid reminder type. reminder:%#v", r)
	}
	return false, nil
}