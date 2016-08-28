package app

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/appengine/lib/linebot"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/appengine/model"
	"github.com/utahta/momoclo-channel/twitter"
	"github.com/utahta/momoclo-channel/ustream"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"
)

type UstreamNotification struct {
	context context.Context
	log     log.Logger
}

func newUstreamNotification(ctx context.Context) *UstreamNotification {
	return &UstreamNotification{context: ctx, log: log.NewGaeLogger(ctx)}
}

func (u *UstreamNotification) Notify() *Error {
	ctx, cancel := context.WithTimeout(u.context, 50*time.Second)
	defer cancel()

	c := ustream.NewClient()
	c.HttpClient.Transport = &urlfetch.Transport{Context: ctx}

	isLive, err := c.IsLive()
	if err != nil {
		return newError(errors.Wrap(err, "Failed to get ustream status"), http.StatusInternalServerError)
	}

	status := model.NewUstreamStatus()
	if err := status.Get(ctx); err != nil && err != datastore.ErrNoSuchEntity {
		return newError(errors.Wrap(err, "Failed to get ustream status from datastore"), http.StatusInternalServerError)
	}

	if status.IsLive == isLive {
		return nil
	}
	status.IsLive = isLive
	status.Put(u.context)

	if isLive {
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()

			tw := twitter.NewMessageClient(
				os.Getenv("TWITTER_CONSUMER_KEY"),
				os.Getenv("TWITTER_CONSUMER_SECRET"),
				os.Getenv("TWITTER_ACCESS_TOKEN"),
				os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
			)
			tw.Log = u.log

			jst := time.FixedZone("Asia/Tokyo", 9*60*60)
			t := time.Now().In(jst)
			if err := tw.Tweet(fmt.Sprintf("momocloTV が配信を開始しました\n%s", t.Format("from 2006/01/02 15:04:05"))); err != nil {
				u.log.Error(err)
				return
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()

			bot, err := linebot.Dial(ctx)
			if err != nil {
				u.log.Error(err)
				return
			}
			defer bot.Close()

			if err := bot.NotifyUstream(); err != nil {
				u.log.Error(err)
				return
			}
		}()

		wg.Wait()
	}
	return nil
}
