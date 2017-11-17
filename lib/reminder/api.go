package reminder

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/adapter/gateway/api/twitter"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/domain/reminder"
	"github.com/utahta/momoclo-channel/infrastructure/log"
	"github.com/utahta/momoclo-channel/lib/config"
	"github.com/utahta/momoclo-channel/lib/linenotify"
	"golang.org/x/sync/errgroup"
)

func Notify(ctx context.Context) error {
	rows, err := reminder.Repository.GetAll(ctx)
	if err != nil {
		return err
	}
	logger := log.NewAppengineLogger(ctx)
	tweeter := twitter.NewTweeter(ctx)

	now := time.Now().In(config.JST)
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
			if _, err := tweeter.Tweet(model.TweetRequest{Text: row.Text}); err != nil {
				logger.Error(ctx, err)
				return err
			}
			return nil
		})

		eg.Go(func() error {
			if err := linenotify.NotifyMessage(ctx, row.Text); err != nil {
				logger.Error(ctx, err)
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
