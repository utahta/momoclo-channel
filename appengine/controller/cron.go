package controller

import (
	"net/http"

	"github.com/utahta/momoclo-channel/appengine/lib/crawler"
	"github.com/utahta/momoclo-channel/appengine/lib/reminder"
	"github.com/utahta/momoclo-channel/appengine/lib/ustream"
)

// Notify reminder
func CronReminder(w http.ResponseWriter, req *http.Request) {
	ctx := getContext(req)

	if err := reminder.Notify(ctx); err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}
}

// Notify ustream
func CronUstream(w http.ResponseWriter, req *http.Request) {
	ctx := getContext(req)

	if err := ustream.Notify(ctx); err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}
}

// Crawling
func CronCrawl(w http.ResponseWriter, req *http.Request) {
	ctx := getContext(req)

	if err := crawler.Crawl(ctx); err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}
}
