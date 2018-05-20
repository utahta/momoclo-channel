package api

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/fukata/golang-stats-api-handler"
	"github.com/go-chi/chi"
	"github.com/utahta/go-linenotify/auth"
	"github.com/utahta/momoclo-channel/api/middleware"
	"github.com/utahta/momoclo-channel/config"
	"github.com/utahta/momoclo-channel/crawler"
	"github.com/utahta/momoclo-channel/customsearch"
	"github.com/utahta/momoclo-channel/dao"
	"github.com/utahta/momoclo-channel/entity"
	"github.com/utahta/momoclo-channel/event"
	"github.com/utahta/momoclo-channel/linebot"
	"github.com/utahta/momoclo-channel/linenotify"
	"github.com/utahta/momoclo-channel/log"
	"github.com/utahta/momoclo-channel/usecase"
	"github.com/utahta/momoclo-channel/ustream"
)

type (
	backendServer struct {
		logger           log.Logger
		transactor       dao.Transactor
		taskQueue        event.TaskQueue
		ustChecker       ustream.StatusChecker
		feedFetcher      crawler.FeedFetcher
		linebotClient    linebot.Client
		imageSearcher    customsearch.ImageSearcher
		linenotifyToken  linenotify.Token
		linenotifyClient linenotify.Client

		reminderRepo         entity.ReminderRepository
		ustreamStatusRepo    entity.UstreamStatusRepository
		latestEntryRepo      entity.LatestEntryRepository
		tweetItemRepo        entity.TweetItemRepository
		lineItemRepo         entity.LineItemRepository
		lineNotificationRepo entity.LineNotificationRepository
	}
)

// NewBackendServer returns backend server.
func NewBackendServer() Server {
	dh := dao.NewDatastoreHandler()
	return &backendServer{
		logger:           log.NewAELogger(),
		transactor:       dao.NewDatastoreTransactor(),
		taskQueue:        event.NewTaskQueue(),
		ustChecker:       ustream.NewStatusChecker(),
		feedFetcher:      crawler.New(),
		linebotClient:    linebot.New(),
		imageSearcher:    customsearch.NewImageSearcher(),
		linenotifyToken:  linenotify.NewToken(),
		linenotifyClient: linenotify.New(),

		reminderRepo:         entity.NewReminderRepository(dh),
		ustreamStatusRepo:    entity.NewUstreamStatusRepository(dh),
		latestEntryRepo:      entity.NewLatestEntryRepository(dh),
		tweetItemRepo:        entity.NewTweetItemRepository(dh),
		lineItemRepo:         entity.NewLineItemRepository(dh),
		lineNotificationRepo: entity.NewLineNotificationRepository(dh),
	}
}

func (s *backendServer) Handle() {
	r := chi.NewRouter()
	r.Use(middleware.AEContext)

	r.Route("/cron", func(r chi.Router) {
		r.Get("/crawl", s.cronCrawl)
		r.Get("/ustream", s.cronUstream)
		r.Get("/reminder", s.cronReminder)
	})

	r.Route("/enqueue", func(r chi.Router) {
		r.Post("/tweets", s.enqueueTweets)
		r.Post("/lines", s.enqueueLines)
	})

	r.Route("/line", func(r chi.Router) {
		r.Route("/bot", func(r chi.Router) {
			r.Post("/callback", s.lineBotCallback)
			r.Get("/help", s.lineBotHelp)
			r.Get("/about", s.lineBotAbout)
		})

		r.Route("/notify", func(r chi.Router) {
			r.Get("/on", s.lineNotifyOn)
			r.Get("/off", s.lineNotifyOff)
			r.HandleFunc("/callback", s.lineNotifyCallback)

			r.Post("/broadcast", s.lineNotifyBroadcast)
			r.Post("/", s.lineNotify)
		})
	})

	r.HandleFunc("/api/stats", stats_api.Handler)

	http.Handle("/", r)
}

// cronReminder checks reminder
func (s *backendServer) cronReminder(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	remind := usecase.NewRemind(s.logger, s.taskQueue, s.reminderRepo)
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
		s.ustChecker,
		s.ustreamStatusRepo,
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
		s.feedFetcher,
		s.taskQueue,
		s.latestEntryRepo,
	)
	crawlFeeds := usecase.NewCrawlFeeds(s.logger, crawlFeed)
	if err := crawlFeeds.Do(ctx); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// enqueueTweets enqueue tweets event
func (s *backendServer) enqueueTweets(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 540*time.Second)
	defer cancel()

	if err := req.ParseForm(); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	item := crawler.FeedItem{}
	if err := event.ParseTask(req.Form, &item); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	enqueueTweets := usecase.NewEnqueueTweets(
		s.logger,
		s.taskQueue,
		s.transactor,
		s.tweetItemRepo,
	)
	params := usecase.EnqueueTweetsParams{FeedItem: item}
	if err := enqueueTweets.Do(ctx, params); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// enqueueLines enqueue lines event
func (s *backendServer) enqueueLines(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 540*time.Second)
	defer cancel()

	if err := req.ParseForm(); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	item := crawler.FeedItem{}
	if err := event.ParseTask(req.Form, &item); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	enqueueLines := usecase.NewEnqueueLines(
		s.logger,
		s.taskQueue,
		s.transactor,
		s.lineItemRepo,
	)
	params := usecase.EnqueueLinesParams{FeedItem: item}
	if err := enqueueLines.Do(ctx, params); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// lineBotCallback handler
func (s *backendServer) lineBotCallback(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	events, err := linebot.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			failResponse(ctx, w, err, http.StatusBadRequest)
			return
		}
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	handleLineBotEvents := usecase.NewHandleLineBotEvents(
		s.logger,
		s.linebotClient,
		s.imageSearcher,
	)
	params := usecase.HandleLineBotEventsParams{Events: events}
	if err := handleLineBotEvents.Do(ctx, params); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// lineBotHelp handler
func (s *backendServer) lineBotHelp(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	tpl := template.Must(template.ParseFiles("public/templates/linebot/help.html"))
	if err := tpl.Execute(w, nil); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// lineBotAbout handler
func (s *backendServer) lineBotAbout(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	tpl := template.Must(template.ParseFiles("public/templates/linebot/about.html"))
	if err := tpl.Execute(w, nil); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// lineNotifyOn redirect to LINE Notify connection page
func (s *backendServer) lineNotifyOn(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	c, err := auth.New(config.C().LineNotify.ClientID, linenotify.CallbackURL())
	if err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "state", Value: c.State, Expires: time.Now().Add(300 * time.Second), Secure: true})

	s.logger.Info(ctx, "Redirect to LINE Notify connection page")

	err = c.Redirect(w, req)
	if err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// lineNotifyOff redirect to LINE Notify revoking page
func (s *backendServer) lineNotifyOff(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	s.logger.Info(ctx, "Redirect to LINE Notify revoking page")

	// official url
	http.Redirect(w, req, "https://notify-bot.line.me/my/", http.StatusFound)
}

// lineNotifyCallback stores LINE Notify token
func (s *backendServer) lineNotifyCallback(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	params, err := auth.ParseRequest(req)
	if err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	state, err := req.Cookie("state")
	if err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	if params.State != state.Value {
		failResponse(ctx, w, errors.New("invalid csrf token"), http.StatusBadRequest)
		return
	}

	addLineNotification := usecase.NewAddLineNotification(
		s.logger,
		s.linenotifyToken,
		s.lineNotificationRepo,
	)
	if err := addLineNotification.Do(ctx, usecase.AddLineNotificationParams{Code: params.Code}); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	t, err := template.New("callback").Parse("<html><body><h1>通知ノフ設定オンにしました（・Θ・）</h1></body></html>")
	if err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// lineNotifyBroadcast invokes broadcast line notification event
func (s *backendServer) lineNotifyBroadcast(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	var messages []linenotify.Message
	if err := event.ParseTask(req.Form, &messages); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	lineNotifyBroadcast := usecase.NewLineNotifyBroadcast(
		s.logger,
		s.taskQueue,
		s.lineNotificationRepo,
	)
	params := usecase.LineNotifyBroadcastParams{Messages: messages}
	if err := lineNotifyBroadcast.Do(ctx, params); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// lineNotify notify users of messages
func (s *backendServer) lineNotify(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
	defer cancel()

	if err := req.ParseForm(); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	var request linenotify.Request
	if err := event.ParseTask(req.Form, &request); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	lineNotify := usecase.NewLineNotify(
		s.logger,
		s.taskQueue,
		s.linenotifyClient,
		s.lineNotificationRepo,
	)
	params := usecase.LineNotifyParams{Request: request}
	if err := lineNotify.Do(ctx, params); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
}
