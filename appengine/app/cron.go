package app

import (
	"net/http"

	"google.golang.org/appengine"
)

type CronHandler struct {
}

func (h *CronHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	var err *Error

	switch r.URL.Path {
	case "/cron/crawl":
		err = newCrawler(ctx).Crawl()
	case "/cron/ustream":
		err = newUstreamNotification(ctx).Notify()
	default:
		http.NotFound(w, r)
	}
	err.Handle(ctx, w)
}
