package app

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/appengine/lib/linebot"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/appengine/model"
	"github.com/utahta/momoclo-channel/crawler"
	"github.com/utahta/momoclo-channel/twitter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

// Queue for crawler.Channel
type QueueHandler struct {
	log log.Logger
}

func (h *QueueHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	h.log = log.NewGaeLogger(ctx)
	var err *Error

	ch, err := h.parseParams(r)
	if err != nil {
		err.Handle(ctx, w)
		return
	}

	switch r.URL.Path {
	case "/queue/tweet":
		err = h.tweet(ctx, ch)
	case "/queue/line":
		err = h.line(ctx, ch)
	default:
		http.NotFound(w, r)
	}
	err.Handle(ctx, w)
}

func (h *QueueHandler) parseParams(r *http.Request) (*crawler.Channel, *Error) {
	var ch crawler.Channel
	if err := json.Unmarshal([]byte(r.FormValue("channel")), &ch); err != nil {
		return nil, newError(errors.Wrapf(err, "Failed to unmarshal."), http.StatusInternalServerError)
	}
	return &ch, nil
}

func (h *QueueHandler) tweet(ctx context.Context, ch *crawler.Channel) *Error {
	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()
	client := urlfetch.Client(ctx)

	var wg sync.WaitGroup
	wg.Add(len(ch.Items))
	for _, item := range ch.Items {
		go func(ctx context.Context, item *crawler.ChannelItem) {
			defer wg.Done()

			if err := model.NewTweetItem(item).Put(ctx); err != nil {
				return
			}

			tw := twitter.NewChannelClient(
				os.Getenv("TWITTER_CONSUMER_KEY"),
				os.Getenv("TWITTER_CONSUMER_SECRET"),
				os.Getenv("TWITTER_ACCESS_TOKEN"),
				os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
			)
			tw.Log = log.NewGaeLogger(ctx)
			tw.Api.HttpClient = client

			if err := tw.TweetItem(ch.Title, item); err != nil {
				h.log.Error(err)
				return
			}
		}(ctx, item)
	}
	wg.Wait()

	return nil
}

func (h *QueueHandler) line(ctx context.Context, ch *crawler.Channel) *Error {
	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(ch.Items))
	for _, item := range ch.Items {
		go func(ctx context.Context, item *crawler.ChannelItem) {
			defer wg.Done()

			if err := model.NewLineItem(item).Put(ctx); err != nil {
				return
			}

			bot, err := linebot.Dial(ctx)
			if err != nil {
				h.log.Error(err)
				return
			}
			defer bot.Close()

			if err := bot.NotifyChannel(ch.Title, item); err != nil {
				h.log.Error(err)
				return
			}
		}(ctx, item)
	}
	wg.Wait()

	return nil
}
