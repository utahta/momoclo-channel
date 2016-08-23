package app

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/log"
	"github.com/utahta/momoclo-channel/model"
	"github.com/utahta/momoclo-channel/twitter"
	"github.com/utahta/momoclo-channel/ustream"
	"golang.org/x/net/context"
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
	c := ustream.NewClient()
	c.HttpClient.Transport = &urlfetch.Transport{Context: u.context}

	isLive, err := c.IsLive()
	if err != nil {
		return newError(errors.Wrap(err, "Failed to get ustream status"), http.StatusInternalServerError)
	}

	status := model.NewUstreamStatus()
	status.Get(u.context)

	if status.IsLive == isLive {
		return nil
	}
	status.IsLive = isLive
	status.Put(u.context)

	if isLive {
		tw := twitter.NewMessageClient(
			os.Getenv("TWITTER_CONSUMER_KEY"),
			os.Getenv("TWITTER_CONSUMER_SECRET"),
			os.Getenv("TWITTER_ACCESS_TOKEN"),
			os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
		)
		tw.Log = u.log

		jst := time.FixedZone("Asia/Tokyo", 9*60*60)
		t := time.Now().In(jst)
		tw.Tweet(fmt.Sprintf("＿人人人人人人人人人人人人人人人人人＿\n＞　momocloTV が配信を開始しました　＜\n￣Y^Y^Y^Y^Y^Y^Y^Y^Y^Y^Y^Y^Y^Y^Y^￣\n%s", t.Format("[2006/01/02 15:04:05]")))
	}
	return nil
}
