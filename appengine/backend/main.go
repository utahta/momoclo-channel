package backend

import (
	"net/http"

	"github.com/fukata/golang-stats-api-handler"
	"github.com/go-chi/chi"
	"github.com/utahta/momoclo-channel/adapter/handler/backend"
	"github.com/utahta/momoclo-channel/adapter/handler/middleware"
	"github.com/utahta/momoclo-channel/lib/config"
)

func init() {
	config.MustLoad("config/deploy.toml")

	router := chi.NewRouter()
	router.Use(middleware.AEContext)

	router.Route("/cron", func(r chi.Router) {
		r.Get("/crawl", backend.CronCrawl)
		r.Get("/ustream", backend.CronUstream)
		r.Get("/reminder", backend.CronReminder)
	})

	router.Route("/enqueue", func(r chi.Router) {
		r.Post("/tweets", backend.EnqueueTweets)
		r.Post("/lines", backend.EnqueueLines)
	})

	router.Route("/line", func(r chi.Router) {
		r.Route("/bot", func(r chi.Router) {
			r.Post("/callback", backend.LineBotCallback)
			r.Get("/help", backend.LineBotHelp)
			r.Get("/about", backend.LineBotAbout)
		})

		r.Route("/notify", func(r chi.Router) {
			r.HandleFunc("/callback", backend.LineNotifyCallback)
			r.Get("/on", backend.LineNotifyOn)
			r.Get("/off", backend.LineNotifyOff)

			r.Post("/broadcast", backend.LineNotifyBroadcast)
			r.Post("/", backend.LineNotify)
		})
	})

	router.HandleFunc("/api/stats", stats_api.Handler)

	http.Handle("/", router)
}
