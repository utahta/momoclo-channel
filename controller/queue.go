package controller

import (
	"net/http"

	"github.com/utahta/momoclo-channel/app"
	"github.com/utahta/momoclo-channel/lib/crawler"
	"github.com/utahta/momoclo-channel/lib/linenotify"
	"github.com/utahta/momoclo-channel/lib/twitter"
)

func QueueTweet(w http.ResponseWriter, req *http.Request) {
	ctx := app.GetContext(req)

	if err := req.ParseForm(); err != nil {
		ctx.Fail(err)
		return
	}

	q := crawler.NewQueueTask()
	ch, err := q.ParseURLValues(req.Form)
	if err != nil {
		ctx.Fail(err)
		return
	}

	if err := twitter.TweetChannel(ctx, ch.(*twitter.ChannelParam)); err != nil {
		ctx.Fail(err)
		return
	}
}

func QueueLine(w http.ResponseWriter, req *http.Request) {
	ctx := app.GetContext(req)

	if err := req.ParseForm(); err != nil {
		ctx.Fail(err)
		return
	}

	q := crawler.NewQueueTask()
	ch, err := q.ParseURLValues(req.Form)
	if err != nil {
		ctx.Fail(err)
		return
	}

	if err := linenotify.NotifyChannel(ctx, ch.(*linenotify.ChannelParam)); err != nil {
		ctx.Fail(err)
		return
	}
}
