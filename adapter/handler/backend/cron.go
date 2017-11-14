package backend

import (
	"context"
	"net/http"
	"time"

	"github.com/utahta/momoclo-channel/adapter/handler"
	"github.com/utahta/momoclo-channel/container"
	"github.com/utahta/momoclo-channel/lib/reminder"
	"github.com/utahta/momoclo-channel/lib/ustream"
)

// CronReminder checks reminder
func CronReminder(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	if err := reminder.Notify(ctx); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// CronUstream checks ustream status
func CronUstream(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	if err := ustream.Notify(ctx); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// CronCrawl crawling some sites
func CronCrawl(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 50*time.Second)
	defer cancel()

	crawl := container.Usecase(ctx).CrawlAll()
	if err := crawl.Do(); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}
}
