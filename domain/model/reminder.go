package model

import (
	"time"

	"github.com/pkg/errors"
)

type (
	// ReminderType type
	ReminderType int

	// Reminder represents any remind information
	Reminder struct {
		ID      int64 `datastore:"-" goon:"id"`
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

	// ReminderRepository interface
	ReminderRepository interface {
		FindAll() ([]*Reminder, error)
		Save(*Reminder) error
	}
)

const (
	ReminderOneTime ReminderType = iota
	ReminderWeekly
)

// NewReminderOnce builds one-time reminder
func NewReminderOnce(text string, remindAt time.Time) *Reminder {
	return &Reminder{
		Text:     text,
		Type:     ReminderOneTime,
		Enabled:  true,
		RemindAt: remindAt,
	}
}

// NewReminderWeekly builds weekly reminder
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

// Disable disabled remind
func (r *Reminder) Disable() {
	r.Enabled = false
}

// IsOneTime returns true if reminder type is one-time
func (r *Reminder) IsOneTime() bool {
	return r.Type == ReminderOneTime
}

// Valid returns true if reminder on time
func (r *Reminder) Valid(now time.Time) (bool, error) {
	switch r.Type {
	case ReminderOneTime:
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
		return false, errors.Errorf("invalid reminder type")
	}
	return false, nil
}
