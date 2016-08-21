package app

import (
	"net/http"
)

type CronHandler struct {
}

func (h *CronHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/cron/crawl":
		new(CrawlHandler).ServeHTTP(w, r)
	default:
		http.NotFound(w, r)
	}
}
