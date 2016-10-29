package reminder

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/go-atomicbool"
	"github.com/utahta/momoclo-channel/appengine/lib/linenotify"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/appengine/lib/twitter"
	"github.com/utahta/momoclo-channel/appengine/model"
	"golang.org/x/net/context"
)

func Notify(ctx context.Context) error {
	q := model.NewReminderQuery(ctx)
	rows, err := q.GetAll()
	if err != nil {
		return err
	}

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now().In(jst)
	for _, row := range rows {
		if ok, err := row.Valid(now); !ok {
			if err != nil {
				return err
			}
			continue
		}

		// Tweet, Line の出し分けが今のところ出来ないので要検討
		const maxGoroutineNum = 2
		errFlg := atomicbool.New(false)
		var wg sync.WaitGroup
		wg.Add(maxGoroutineNum)

		go func(text string) {
			defer wg.Done()
			if err := twitter.TweetMessage(ctx, text); err != nil {
				errFlg.Set(true)
				log.GaeLog(ctx).Error(err)
			}
		}(row.Text)

		go func(text string) {
			defer wg.Done()
			if err := linenotify.NotifyMessage(ctx, text); err != nil {
				errFlg.Set(true)
				log.GaeLog(ctx).Error(err)
			}
		}(row.Text)

		wg.Wait()

		if row.IsOnce() {
			row.Disable(ctx)
		}

		if errFlg.Enabled() {
			return errors.Errorf("Errors occured in reminder.Notify. text:%s", row.Text)
		}
	}
	return nil
}
