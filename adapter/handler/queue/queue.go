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

	var requests []model.TweetRequest
	if err := event.ParseTask(req.Form, &requests); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	params := usecase.TweetParams{TweetRequests: requests}
	if err := container.Usecase(ctx).Tweet().Do(params); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// LineNotifyBroadcast invokes broadcast line notification event
func LineNotifyBroadcast(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	var messages []model.LineNotifyMessage
	if err := event.ParseTask(req.Form, &messages); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	params := usecase.LineNotifyBroadcastParams{Messages: messages}
	if err := container.Usecase(ctx).LineNotifyBroadcast().Do(params); err != nil {
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

	var request model.LineNotifyRequest
	if err := event.ParseTask(req.Form, &request); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	params := usecase.LineNotifyParams{Request: request}
	if err := container.Usecase(ctx).LineNotify().Do(params); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}
}
