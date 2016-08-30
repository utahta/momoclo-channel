package app

import (
	"net/http"
	"sync"
	"time"

	"github.com/utahta/momoclo-channel/appengine/lib/linebot"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/appengine/model"
	"golang.org/x/net/context"
)

type ReminderNotification struct {
	context context.Context
	log     log.Logger
}

func newReminderNotification(ctx context.Context) *ReminderNotification {
	return &ReminderNotification{context: ctx, log: log.NewGaeLogger(ctx)}
}

func (r *ReminderNotification) Notify() *Error {
	ctx, cancel := context.WithTimeout(r.context, 50*time.Second)
	defer cancel()

	q := model.NewReminderQuery(ctx)
	rows, err := q.GetAll()
	if err != nil {
		return newError(err, http.StatusInternalServerError)
	}

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now().In(jst)
	for _, row := range rows {
		if ok, err := row.Valid(now); !ok || err != nil {
			r.log.Error(err)
			continue
		}

		var wg sync.WaitGroup

		//wg.Add(1)
		//go func() {
		//	defer wg.Done()
		//}()
		//
		wg.Add(1)
		go func(text string) {
			defer wg.Done()
			linebot.NotifyReminder(ctx, text)
		}(row.Text)

		wg.Wait()

		if row.IsOnce() {
			row.Disable(ctx)
		}
	}
	return nil
}
