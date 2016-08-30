package app

import (
	"net/http"
	"time"

	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

type LinebotHandler struct {
	log log.Logger
}

func (h *LinebotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	h.log = log.NewGaeLogger(ctx)
	var err *Error

	switch r.URL.Path {
	case "/linebot/callback":
		err = h.callback(ctx)
	default:
		http.NotFound(w, r)
	}
	err.Handle(ctx, w)
}

func (h *LinebotHandler) callback(ctx context.Context) *Error {
	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()
	
	return nil
}
