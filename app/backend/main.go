package backend

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/utahta/momoclo-channel/controller"
	"github.com/utahta/momoclo-channel/lib/config"
	"github.com/utahta/momoclo-channel/middleware"
)

func init() {
	config.MustLoad("config/deploy.toml")

	router := chi.NewRouter()
	router.Use(middleware.Appengine)

	router.Get("/cron/crawl", controller.CronCrawl)
	router.Get("/cron/ustream", controller.CronUstream)
	router.Get("/cron/reminder", controller.CronReminder)

	router.Post("/linebot/callback", controller.LineBotCallback)
	router.Get("/linebot/help", controller.LineBotHelp)
	router.Get("/linebot/about", controller.LineBotAbout)

	router.HandleFunc("/linenotify/callback", controller.LinenotifyCallback)
	router.Get("/linenotify/on", controller.LinenotifyOn)
	router.Get("/linenotify/off", controller.LinenotifyOff)

	http.Handle("/", router)

	log.Println("init backend app")
}
