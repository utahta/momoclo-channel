package reminder

import (
	"sync"
	"time"

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
				log.GaeLog(ctx).Error(err)
			}
			continue
		}

		var wg sync.WaitGroup

		wg.Add(1)
		go func(text string) {
			defer wg.Done()
			twitter.TweetText(ctx, text)
		}(row.Text)

		wg.Add(1)
		go func(text string) {
			defer wg.Done()
			linenotify.NotifyMessage(ctx, text)
		}(row.Text)

		wg.Wait()

		if row.IsOnce() {
			row.Disable(ctx)
		}
	}
	return nil
}
