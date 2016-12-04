package controller

import (
	"net/http"

	"github.com/utahta/momoclo-channel/appengine/lib/crawler"
	"github.com/utahta/momoclo-channel/appengine/lib/linenotify"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/appengine/lib/twitter"
)

func QueueTweet(w http.ResponseWriter, req *http.Request) {
	ctx := getContext(req)

	if err := req.ParseForm(); err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}

	q := crawler.NewQueueTask(log.NewGaeLogger(ctx))
	ch, err := q.ParseURLValues(req.Form)
	if err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}

	if err := twitter.TweetChannel(ctx, ch); err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}
}

func QueueLine(w http.ResponseWriter, req *http.Request) {
	ctx := getContext(req)

	if err := req.ParseForm(); err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}

	q := crawler.NewQueueTask(log.NewGaeLogger(ctx))
	ch, err := q.ParseURLValues(req.Form)
	if err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}

	if err := linenotify.NotifyChannel(ctx, ch); err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}
}
