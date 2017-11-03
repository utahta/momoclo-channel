package queue

import (
	"net/http"

	"github.com/utahta/momoclo-channel/adapter/handler"
	"github.com/utahta/momoclo-channel/lib/crawler"
	"github.com/utahta/momoclo-channel/lib/linenotify"
	"github.com/utahta/momoclo-channel/lib/twitter"
)

func QueueTweet(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	param := &twitter.ChannelParam{}
	q := crawler.NewQueueTask()
	if err := q.ParseURLValues(req.Form, param); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	if err := twitter.TweetChannel(ctx, param); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

func QueueLine(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	param := &linenotify.ChannelParam{}
	q := crawler.NewQueueTask()
	if err := q.ParseURLValues(req.Form, param); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	if err := linenotify.NotifyChannel(ctx, param); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}
}
