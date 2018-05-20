package api

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/utahta/momoclo-channel/api/middleware"
	"github.com/utahta/momoclo-channel/event"
	"github.com/utahta/momoclo-channel/log"
	"github.com/utahta/momoclo-channel/twitter"
	"github.com/utahta/momoclo-channel/usecase"
)

type (
	batchServer struct {
		logger    log.Logger
		taskQueue event.TaskQueue
		tweeter   twitter.Tweeter
	}
)

// NewBatchServer returns batch server.
func NewBatchServer() Server {
	return &batchServer{
		logger:    log.NewAELogger(),
		taskQueue: event.NewTaskQueue(),
		tweeter:   twitter.NewTweeter(),
	}
}

func (s *batchServer) Handle() {
	r := chi.NewRouter()
	r.Use(middleware.AEContext)

	r.Get("/_ah/start", func(w http.ResponseWriter, req *http.Request) {}) // nop
	r.Post("/tweet", s.tweet)

	http.Handle("/", r)
}

// tweet tweets messages
func (s *batchServer) tweet(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 180*time.Second) //TODO ctx should be an argument?
	defer cancel()

	if err := req.ParseForm(); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	var requests []twitter.TweetRequest
	if err := event.ParseTask(req.Form, &requests); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	tweet := usecase.NewTweet(
		s.logger,
		s.taskQueue,
		s.tweeter,
	)
	params := usecase.TweetParams{Requests: requests}
	if err := tweet.Do(ctx, params); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
}
