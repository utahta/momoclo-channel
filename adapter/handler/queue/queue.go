package queue

import (
	"net/http"

	"github.com/utahta/momoclo-channel/adapter/handler"
	"github.com/utahta/momoclo-channel/container"
	"github.com/utahta/momoclo-channel/lib/linenotify"
	"github.com/utahta/momoclo-channel/lib/twitter"
)

// Tweet invokes tweet event
func Tweet(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	param := &twitter.ChannelParam{}
	crawl := container.Usecase(ctx).Crawl()
	if err := crawl.ParseURLValues(req.Form, param); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	if err := twitter.TweetChannel(ctx, param); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// LineNotify invokes line notification event
func LineNotify(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	param := &linenotify.ChannelParam{}
	crawl := container.Usecase(ctx).Crawl()
	if err := crawl.ParseURLValues(req.Form, param); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	if err := linenotify.NotifyChannel(ctx, param); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}
}
