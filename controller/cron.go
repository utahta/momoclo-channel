package controller

import (
	"net/http"

	"github.com/utahta/momoclo-channel/lib/crawler"
	"github.com/utahta/momoclo-channel/lib/reminder"
	"github.com/utahta/momoclo-channel/lib/ustream"
)

// Notify reminder
func CronReminder(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	if err := reminder.Notify(ctx); err != nil {
		fail(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// Notify ustream
func CronUstream(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	if err := ustream.Notify(ctx); err != nil {
		fail(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// Crawling
func CronCrawl(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	if err := crawler.Crawl(ctx); err != nil {
		fail(ctx, w, err, http.StatusInternalServerError)
		return
	}
}
