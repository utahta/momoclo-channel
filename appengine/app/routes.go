package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"github.com/utahta/momoclo-channel/appengine/controller"
	"github.com/utahta/momoclo-channel/appengine/middleware"
)

func initRoutes() {
	n := negroni.New()
	n.Use(negroni.HandlerFunc(middleware.Appengine))

	router := mux.NewRouter()
	router.HandleFunc("/cron/crawl", controller.CronCrawl).Methods("GET")
	router.HandleFunc("/cron/ustream", controller.CronUstream).Methods("GET")
	router.HandleFunc("/cron/reminder", controller.CronReminder).Methods("GET")

	router.HandleFunc("/queue/tweet", controller.QueueTweet).Methods("POST")
	router.HandleFunc("/queue/line", controller.QueueLine).Methods("POST")

	router.HandleFunc("/linebot/callback", controller.LineBotCallback).Methods("POST")
	router.HandleFunc("/linebot/help", controller.LineBotHelp).Methods("GET")
	router.HandleFunc("/linebot/about", controller.LineBotAbout).Methods("GET")

	router.HandleFunc("/linenotify/callback", controller.LinenotifyCallback)
	router.HandleFunc("/linenotify/on", controller.LinenotifyCallback).Methods("GET")
	router.HandleFunc("/linenotify/off", controller.LinenotifyCallback).Methods("GET")

	n.UseHandler(router)
	http.Handle("/", n)
}
