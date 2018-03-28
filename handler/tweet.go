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

// Tweet tweets messages
func Tweet(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 180*time.Second) //TODO ctx should be an argument?
	defer cancel()

	if err := req.ParseForm(); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	var requests []types.TweetRequest
	if err := event.ParseTask(req.Form, &requests); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	params := usecase.TweetParams{Requests: requests}
	if err := container.Usecase(ctx).Tweet().Do(params); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
}
