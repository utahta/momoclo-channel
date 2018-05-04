package backend

import (
	"net/http"

	"github.com/fukata/golang-stats-api-handler"
	"github.com/go-chi/chi"
	"github.com/utahta/momoclo-channel/config"
	"github.com/utahta/momoclo-channel/handler"
	"github.com/utahta/momoclo-channel/handler/middleware"
)

func init() {
	config.MustLoad("config/deploy.toml")

	router := chi.NewRouter()
	router.Use(middleware.AEContext)

	router.Route("/cron", func(r chi.Router) {
		r.Get("/crawl", handler.CronCrawl)
		r.Get("/ustream", handler.CronUstream)
		r.Get("/reminder", handler.CronReminder)
	})

	router.Route("/enqueue", func(r chi.Router) {
		r.Post("/tweets", handler.EnqueueTweets)
		r.Post("/lines", handler.EnqueueLines)
	})

	router.Route("/line", func(r chi.Router) {
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

	router.HandleFunc("/api/stats", stats_api.Handler)

	http.Handle("/", router)
}
