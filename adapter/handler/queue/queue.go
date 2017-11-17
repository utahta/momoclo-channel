package queue

import (
	"context"
	"net/http"
	"time"

	"github.com/utahta/momoclo-channel/adapter/handler"
	"github.com/utahta/momoclo-channel/container"
	"github.com/utahta/momoclo-channel/domain/event"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/usecase"
)

// Tweet tweets to twitter
func Tweet(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 540*time.Second)
	defer cancel()

	if err := req.ParseForm(); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	var tweetRequests []model.TweetRequest
	if err := event.ParseTask(req.Form, &tweetRequests); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	params := usecase.TweetParams{TweetRequests: tweetRequests}
	if err := container.Usecase(ctx).Tweet().Do(params); err != nil {
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

	//param := &linenotify.ChannelParam{}
	//crawl := container.Usecase(ctx).CrawlAll()
	//if err := crawl.ParseURLValues(req.Form, param); err != nil {
	//	handler.Fail(ctx, w, err, http.StatusInternalServerError)
	//	return
	//}

	//if err := linenotify.NotifyChannel(ctx, param); err != nil {
	//	handler.Fail(ctx, w, err, http.StatusInternalServerError)
	//	return
	//}
}
