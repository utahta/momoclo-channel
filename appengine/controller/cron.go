package controller

import (
	"net/http"
	"time"

	"github.com/utahta/momoclo-channel/appengine/lib/crawler"
	"github.com/utahta/momoclo-channel/appengine/lib/reminder"
	"github.com/utahta/momoclo-channel/appengine/lib/ustream"
	"golang.org/x/net/context"
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
	ctx, cancel := context.WithTimeout(getContext(req), 30*time.Second)
	defer cancel()

	if err := crawler.Crawl(ctx); err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}
}
