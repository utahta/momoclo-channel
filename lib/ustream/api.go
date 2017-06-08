package ustream

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/go-atomicbool"
	"github.com/utahta/momoclo-channel/lib/linenotify"
	"github.com/utahta/momoclo-channel/lib/log"
	"github.com/utahta/momoclo-channel/lib/twitter"
	"github.com/utahta/momoclo-channel/model"
	"github.com/utahta/uststat"
	"golang.org/x/net/context"
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
		const maxGoroutineNum = 2
		errFlg := atomicbool.New(false)
		var wg sync.WaitGroup
		wg.Add(maxGoroutineNum)

		go func() {
			defer wg.Done()
			jst := time.FixedZone("Asia/Tokyo", 9*60*60)
			t := time.Now().In(jst)
			if err := twitter.TweetMessage(ctx, fmt.Sprintf("momocloTV が配信を開始しました\n%s\nhttp://www.ustream.tv/channel/momoclotv", t.Format("from 2006/01/02 15:04:05"))); err != nil {
				errFlg.Set(true)
				log.Error(ctx, err)
			}
		}()

		go func() {
			defer wg.Done()
			if err := linenotify.NotifyMessage(ctx, "momocloTV が配信を開始しました\nhttp://www.ustream.tv/channel/momoclotv"); err != nil {
				errFlg.Set(true)
				log.Error(ctx, err)
			}
		}()

		wg.Wait()

		if errFlg.Enabled() {
			return errors.New("Errors occurred in ustream.Notify")
		}
	}
	return nil
}
