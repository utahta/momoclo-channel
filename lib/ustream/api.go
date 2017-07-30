package ustream

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/lib/config"
	"github.com/utahta/momoclo-channel/lib/linenotify"
	"github.com/utahta/momoclo-channel/lib/twitter"
	"github.com/utahta/momoclo-channel/model"
	"github.com/utahta/uststat"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"
)

func Notify(ctx context.Context) error {
	c, err := uststat.New(uststat.WithHTTPTransport(&urlfetch.Transport{Context: ctx}))
	if err != nil {
		return err
	}

	isLive, err := c.IsLiveByChannelID("4979543")
	if err != nil {
		return errors.Wrap(err, "Failed to get ustream status")
	}

	status := model.NewUstreamStatus()
	if err := status.Get(ctx); err != nil && err != datastore.ErrNoSuchEntity {
		return errors.Wrap(err, "Failed to get ustream status from datastore")
	}

	if status.IsLive == isLive {
		return nil
	}
	status.IsLive = isLive
	status.Put(ctx)

	if isLive {
		eg := &errgroup.Group{}

		eg.Go(func() error {
			t := time.Now().In(config.JST)
			if err := twitter.TweetMessage(ctx, fmt.Sprintf("momocloTV が配信を開始しました\n%s\nhttp://www.ustream.tv/channel/momoclotv", t.Format("from 2006/01/02 15:04:05"))); err != nil {
				return err
			}
			return nil
		})

		eg.Go(func() error {
			if err := linenotify.NotifyMessage(ctx, "momocloTV が配信を開始しました\nhttp://www.ustream.tv/channel/momoclotv"); err != nil {
				return err
			}
			return nil
		})

		if err := eg.Wait(); err != nil {
			return errors.Wrap(err, "Errors occurred in ustream.Notify")
		}
	}
	return nil
}
