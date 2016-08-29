package model

import (
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type ReminderOnce struct {
	Id        int64 `datastore:"-" goon:"id"`
	Text      string
	RemindAt  time.Time
	CreatedAt time.Time
}

type ReminderOnceQuery struct {
	context context.Context
}

func NewReminderOnceQuery(ctx context.Context) *ReminderOnceQuery {
	return &ReminderOnceQuery{context: ctx}
}

func (r *ReminderOnceQuery) GetAll() ([]*ReminderOnce, error) {
	t := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
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
