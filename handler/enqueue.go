package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/utahta/momoclo-channel/container"
	"github.com/utahta/momoclo-channel/event"
	"github.com/utahta/momoclo-channel/types"
	"github.com/utahta/momoclo-channel/usecase"
)

// EnqueueTweets enqueue tweets event
func EnqueueTweets(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 540*time.Second)
	defer cancel()

	if err := req.ParseForm(); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	item := types.FeedItem{}
	if err := event.ParseTask(req.Form, &item); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	params := usecase.EnqueueTweetsParams{FeedItem: item}
	if err := container.Usecase(ctx).EnqueueTweets().Do(params); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// EnqueueLines enqueue lines event
func EnqueueLines(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 540*time.Second)
	defer cancel()

	if err := req.ParseForm(); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	item := types.FeedItem{}
	if err := event.ParseTask(req.Form, &item); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	params := usecase.EnqueueLinesParams{FeedItem: item}
	if err := container.Usecase(ctx).EnqueueLines().Do(params); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
}
