package reminder

import (
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/lib/linenotify"
	"github.com/utahta/momoclo-channel/lib/log"
	"github.com/utahta/momoclo-channel/lib/twitter"
	"github.com/utahta/momoclo-channel/model"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
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
		eg := new(errgroup.Group)

		eg.Go(func() error {
			if err := twitter.TweetMessage(ctx, row.Text); err != nil {
				log.Error(ctx, err)
				return err
			}
			return nil
		})
		eg.Go(func() error {
			if err := linenotify.NotifyMessage(ctx, row.Text); err != nil {
				log.Error(ctx, err)
				return err
			}
			return nil
		})

		err := eg.Wait()

		if row.IsOnce() {
			row.Disable(ctx)
		}

		if err != nil {
			return errors.Errorf("Errors occurred in reminder.Notify. text:%s", row.Text)
		}
	}
	return nil
}
