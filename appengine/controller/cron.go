package controller

import (
	"net/http"

	"github.com/utahta/momoclo-channel/appengine/lib/crawler"
	"github.com/utahta/momoclo-channel/appengine/lib/reminder"
	"github.com/utahta/momoclo-channel/appengine/lib/ustream"
	"golang.org/x/net/context"
)

// Notify reminder
func CronReminder(ctx context.Context, w http.ResponseWriter, req *http.Request) *Error {
	if err := reminder.Notify(ctx); err != nil {
		return newError(err, http.StatusInternalServerError)
	}
	return nil
}

// Notify ustream
func CronUstream(ctx context.Context, w http.ResponseWriter, req *http.Request) *Error {
	if err := ustream.Notify(ctx); err != nil {
		return newError(err, http.StatusInternalServerError)
	}
	return nil
}

// Crawling
func CronCrawl(ctx context.Context, w http.ResponseWriter, req *http.Request) *Error {
	if err := crawler.Crawl(ctx); err != nil {
		return newError(err, http.StatusInternalServerError)
	}
	return nil
}
