package controller

import (
	"net/http"

	"github.com/utahta/momoclo-channel/appengine/lib/crawler"
	"github.com/utahta/momoclo-channel/appengine/lib/linenotify"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/appengine/lib/twitter"
	"golang.org/x/net/context"
)

func QueueTweet(ctx context.Context, w http.ResponseWriter, req *http.Request) *Error {
	if err := req.ParseForm(); err != nil {
		return newError(err, http.StatusInternalServerError)
	}

	q := crawler.NewQueueTask(log.NewGaeLogger(ctx))
	ch, err := q.ParseURLValues(req.Form)
	if err != nil {
		return newError(err, http.StatusInternalServerError)
	}

	if err := twitter.TweetChannel(ctx, ch); err != nil {
		return newError(err, http.StatusInternalServerError)
	}
	return nil
}

func QueueLine(ctx context.Context, w http.ResponseWriter, req *http.Request) *Error {
	if err := req.ParseForm(); err != nil {
		return newError(err, http.StatusInternalServerError)
	}

	q := crawler.NewQueueTask(log.NewGaeLogger(ctx))
	ch, err := q.ParseURLValues(req.Form)
	if err != nil {
		return newError(err, http.StatusInternalServerError)
	}

	if err := linenotify.NotifyChannel(ctx, ch); err != nil {
		return newError(err, http.StatusInternalServerError)
	}
	return nil
}
