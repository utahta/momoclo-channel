package app

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
	"sync"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/crawler"
	"github.com/utahta/momoclo-channel/log"
	"github.com/utahta/momoclo-channel/model"
	"github.com/utahta/momoclo-channel/twitter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

type QueueHandler struct {
	context context.Context
}

func (h *QueueHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.context = appengine.NewContext(r)

	var err *Error
	switch r.URL.Path {
	case "/queue/tweet":
		err = h.serveTweet(w, r)
	case "/queue/line":
		err = h.serveLine(w, r)
	default:
		http.NotFound(w, r)
	}

	err.Handle(h.context, w)
}

func (h *QueueHandler) parseParams(r *http.Request) (*crawler.Channel, error) {
	var ch crawler.Channel
	if err := json.Unmarshal([]byte(r.FormValue("channel")), &ch); err != nil {
		return nil, errors.Wrapf(err, "Failed to unmarshal.")
	}
	return &ch, nil
}

func (h *QueueHandler) serveTweet(w http.ResponseWriter, r *http.Request) *Error {
	ch, err := h.parseParams(r)
	if err != nil {
		return newError(err, http.StatusInternalServerError)
	}

	var wg sync.WaitGroup
	for _, item := range ch.Items {
		wg.Add(1)
		go func(item *crawler.ChannelItem) {
			if err := model.PutTweetItem(h.context, item); err != nil {
				return
			}

			ctx, cancel := context.WithTimeout(h.context, 45*time.Second)
			defer func() {
				cancel()
				wg.Done()
			}()

			tw := twitter.NewTwitterClient(
				os.Getenv("TWITTER_CONSUMER_KEY"),
				os.Getenv("TWITTER_CONSUMER_SECRET"),
				os.Getenv("TWITTER_ACCESS_TOKEN"),
				os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
			)
			tw.Log = log.NewGaeLogger(h.context)
			tw.Api.HttpClient.Transport = &urlfetch.Transport{Context: ctx}

			tw.TweetItem(ch.Title, item)
		}(item)
	}
	wg.Wait()

	return nil
}

func (h *QueueHandler) serveLine(w http.ResponseWriter, r *http.Request) *Error {
	var ch crawler.Channel
	if err := json.Unmarshal([]byte(r.FormValue("channel")), &ch); err != nil {
		return newError(errors.Wrapf(err, "Failed to unmarshal."), http.StatusInternalServerError)
	}

	return nil
}
