package app

import (
	"net/http"
	
	"google.golang.org/appengine"
)

type CronHandler struct {
}

func (h *CronHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	switch r.URL.Path {
	case "/cron/crawl":
		new(Crawler).Crawl(ctx)
	default:
		http.NotFound(w, r)
	}
}
