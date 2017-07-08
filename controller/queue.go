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

	param := &twitter.ChannelParam{}
	q := crawler.NewQueueTask()
	if err := q.ParseURLValues(req.Form, param); err != nil {
		ctx.Fail(err)
		return
	}

	if err := twitter.TweetChannel(ctx, param); err != nil {
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

	param := &linenotify.ChannelParam{}
	q := crawler.NewQueueTask()
	if err := q.ParseURLValues(req.Form, param); err != nil {
		ctx.Fail(err)
		return
	}

	if err := linenotify.NotifyChannel(ctx, param); err != nil {
		ctx.Fail(err)
		return
	}
}
