package backend

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/utahta/momoclo-channel/adapter/handler/backend"
	"github.com/utahta/momoclo-channel/adapter/handler/middleware"
	"github.com/utahta/momoclo-channel/lib/config"
)

func init() {
	config.MustLoad("config/deploy.toml")

	router := chi.NewRouter()
	router.Use(middleware.AppengineContext)

	router.Get("/cron/crawl", backend.CronCrawl)
	router.Get("/cron/ustream", backend.CronUstream)
	router.Get("/cron/reminder", backend.CronReminder)

	router.Post("/enqueue/tweets", backend.EnqueueTweets)
	router.Post("/enqueue/lines", backend.EnqueueLines)

	router.Post("/linebot/callback", backend.LineBotCallback)
	router.Get("/linebot/help", backend.LineBotHelp)
	router.Get("/linebot/about", backend.LineBotAbout)

	router.HandleFunc("/linenotify/callback", backend.LinenotifyCallback)
	router.Get("/linenotify/on", backend.LinenotifyOn)
	router.Get("/linenotify/off", backend.LinenotifyOff)

	http.Handle("/", router)
}
