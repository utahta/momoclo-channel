package api

import (
	"context"
	"net/http"
	"time"

	"github.com/fukata/golang-stats-api-handler"
	"github.com/go-chi/chi"
	"github.com/utahta/momoclo-channel/api/middleware"
	"github.com/utahta/momoclo-channel/crawler"
	"github.com/utahta/momoclo-channel/dao"
	"github.com/utahta/momoclo-channel/entity"
	"github.com/utahta/momoclo-channel/event"
	"github.com/utahta/momoclo-channel/handler"
	"github.com/utahta/momoclo-channel/log"
	"github.com/utahta/momoclo-channel/usecase"
	"github.com/utahta/momoclo-channel/ustream"
)

type (
	backendServer struct {
		router     *chi.Mux
		daoHandler dao.PersistenceHandler
		logger     log.Logger
		taskQueue  event.TaskQueue
	}
)

func (s *backendServer) routes() {
	s.router.Use(middleware.AEContext)

	s.router.Route("/cron", func(r chi.Router) {
		r.Get("/crawl", s.cronCrawl)
		r.Get("/ustream", s.cronUstream)
		r.Get("/reminder", s.cronReminder)
	})

	s.router.Route("/enqueue", func(r chi.Router) {
		r.Post("/tweets", handler.EnqueueTweets)
		r.Post("/lines", handler.EnqueueLines)
	})

	s.router.Route("/line", func(r chi.Router) {
		r.Route("/bot", func(r chi.Router) {
			r.Post("/callback", handler.LineBotCallback)
			r.Get("/help", handler.LineBotHelp)
			r.Get("/about", handler.LineBotAbout)
		})

		r.Route("/notify", func(r chi.Router) {
			r.HandleFunc("/callback", handler.LineNotifyCallback)
			r.Get("/on", handler.LineNotifyOn)
			r.Get("/off", handler.LineNotifyOff)

			r.Post("/broadcast", handler.LineNotifyBroadcast)
			r.Post("/", handler.LineNotify)
		})
	})

	s.router.HandleFunc("/api/stats", stats_api.Handler)
}

// cronReminder checks reminder
func (s *backendServer) cronReminder(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	remind := usecase.NewRemind(s.logger, s.taskQueue, entity.NewReminderRepository(s.daoHandler))
	if err := remind.Do(ctx); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// cronUstream checks ustream status
func (s *backendServer) cronUstream(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	checkUstream := usecase.NewCheckUstream(
		s.logger,
		s.taskQueue,
		ustream.NewStatusChecker(),
		entity.NewUstreamStatusRepository(s.daoHandler),
	)
	if err := checkUstream.Do(ctx); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// cronCrawl crawling some sites
func (s *backendServer) cronCrawl(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
	defer cancel()

	crawlFeed := usecase.NewCrawlFeed(
		s.logger,
		crawler.New(),
		s.taskQueue,
		entity.NewLatestEntryRepository(s.daoHandler),
	)
	crawlFeeds := usecase.NewCrawlFeeds(s.logger, crawlFeed)
	if err := crawlFeeds.Do(ctx); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
}
