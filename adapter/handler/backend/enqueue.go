package backend

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

// EnqueueTweets enqueue tweets event
func EnqueueTweets(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 540*time.Second)
	defer cancel()

	if err := req.ParseForm(); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	item := model.FeedItem{}
	if err := event.ParseTask(req.Form, &item); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	params := usecase.EnqueueTweetsParams{FeedItem: item}
	if err := container.Usecase(ctx).EnqueueTweets().Do(params); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// EnqueueLines enqueue lines event
func EnqueueLines(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 540*time.Second)
	defer cancel()

	if err := req.ParseForm(); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	item := model.FeedItem{}
	if err := event.ParseTask(req.Form, &item); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	//FIXME implement
}
