package app

import (
	"net/http"
	"encoding/json"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"github.com/utahta/momoclo-channel/crawler"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
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
	default:
		http.NotFound(w, r)
	}

	err.Handle(h.context, w)
}

func (h *QueueHandler) serveTweet(w http.ResponseWriter, r *http.Request) *Error {
	items := []*crawler.ChannelItem{}
	if err := json.Unmarshal([]byte(r.FormValue("items")), &items); err != nil {
		return newError(errors.Wrapf(err, "Failed to unmarshal."), http.StatusInternalServerError)
	}

	for _, item := range items {
		log.Infof(h.context, "%v", item)
	}
	return nil
}
