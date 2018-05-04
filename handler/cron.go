package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/utahta/momoclo-channel/container"
)

// CronReminder checks reminder
func CronReminder(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	if err := container.Usecase(ctx).Remind().Do(); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// CronUstream checks ustream status
func CronUstream(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	if err := container.Usecase(ctx).CheckUstream().Do(); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// CronCrawl crawling some sites
func CronCrawl(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
	defer cancel()

	if err := container.Usecase(ctx).CrawlFeeds().Do(); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
}
