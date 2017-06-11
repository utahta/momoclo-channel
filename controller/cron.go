package controller

import (
	"net/http"

	"github.com/utahta/momoclo-channel/app"
	"github.com/utahta/momoclo-channel/lib/crawler"
	"github.com/utahta/momoclo-channel/lib/reminder"
	"github.com/utahta/momoclo-channel/lib/ustream"
)

// Notify reminder
func CronReminder(w http.ResponseWriter, req *http.Request) {
	ctx := app.GetContext(req)

	if err := reminder.Notify(ctx); err != nil {
		ctx.Fail(err)
		return
	}
}

// Notify ustream
func CronUstream(w http.ResponseWriter, req *http.Request) {
	ctx := app.GetContext(req)

	if err := ustream.Notify(ctx); err != nil {
		ctx.Fail(err)
		return
	}
}

// Crawling
func CronCrawl(w http.ResponseWriter, req *http.Request) {
	ctx := app.GetContext(req)

	if err := crawler.Crawl(ctx); err != nil {
		ctx.Fail(err)
		return
	}
}
