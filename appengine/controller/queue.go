package controller

import (
	"net/http"

	"github.com/utahta/momoclo-channel/appengine/lib/crawler"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"golang.org/x/net/context"
)

func QueueTweet(ctx context.Context, w http.ResponseWriter, req *http.Request) *Error {
	if err := req.ParseForm(); err != nil {
		return newError(err, http.StatusInternalServerError)
	}

	q := crawler.NewQueueTask(log.NewGaeLogger(ctx))
	if err := q.RunTweet(ctx, req.Form); err != nil {
		return newError(err, http.StatusInternalServerError)
	}
	return nil
}

func QueueLine(ctx context.Context, w http.ResponseWriter, req *http.Request) *Error {
	if err := req.ParseForm(); err != nil {
		return newError(err, http.StatusInternalServerError)
	}

	q := crawler.NewQueueTask(log.NewGaeLogger(ctx))
	if err := q.RunLine(ctx, req.Form); err != nil {
		return newError(err, http.StatusInternalServerError)
	}
	return nil
}
