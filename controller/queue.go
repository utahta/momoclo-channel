package controller

import (
	"net/http"

	"github.com/utahta/momoclo-channel/lib/crawler"
	"github.com/utahta/momoclo-channel/lib/linenotify"
	"github.com/utahta/momoclo-channel/lib/twitter"
)

func QueueTweet(w http.ResponseWriter, req *http.Request) {
	ctx := getContext(req)

	if err := req.ParseForm(); err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}

	q := crawler.NewQueueTask()
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

	q := crawler.NewQueueTask()
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
