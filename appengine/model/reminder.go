package model

import (
	"time"

	"github.com/mjibson/goon"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type ReminderOnce struct {
	Id        int64 `datastore:"-" goon:"id"`
	Message   string
	RemindAt  time.Time
	CreatedAt time.Time
}

func NewReminderOnce() *ReminderOnce {
	return &ReminderOnce{
		Message:  "test",
		RemindAt: time.Now(),
	}
}

func (r *ReminderOnce) Put(ctx context.Context) {
	g := goon.FromContext(ctx)
	g.Put(r)
}

type ReminderOnceQuery struct {
	context context.Context
}

func NewReminderOnceQuery(ctx context.Context) *ReminderOnceQuery {
	return &ReminderOnceQuery{context: ctx}
}

func (r *ReminderOnceQuery) GetAll() ([]*ReminderOnce, error) {
	t := time.Now()
	duration := time.Duration(t.Second())*time.Second + time.Duration(t.Nanosecond())*time.Nanosecond
	t = t.Add(-duration)
	q := datastore.NewQuery("ReminderOnce").Filter("RemindAt >=", t).Order("RemindAt")

	var dst []*ReminderOnce
	_, err := q.GetAll(r.context, &dst)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get all ReminderOnce entities.")
	}
	return dst, nil
}
