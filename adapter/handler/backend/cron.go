package backend

import (
	"context"
	"net/http"
	"time"

	"github.com/utahta/momoclo-channel/adapter/handler"
	"github.com/utahta/momoclo-channel/adapter/persistence"
	"github.com/utahta/momoclo-channel/infrastructure/dao"
	"github.com/utahta/momoclo-channel/lib/reminder"
	"github.com/utahta/momoclo-channel/lib/ustream"
	"github.com/utahta/momoclo-channel/usecase"
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
	ctx := req.Context()
	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	crawl := usecase.NewCrawl(persistence.NewLatestEntryRepository(dao.NewDatastoreHandler(ctx)))
	if err := crawl.Do(ctx); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}
}
