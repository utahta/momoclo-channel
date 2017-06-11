package backend

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/urfave/negroni"
	"github.com/utahta/momoclo-channel/controller"
	"github.com/utahta/momoclo-channel/middleware"
)

func init() {
	if err := godotenv.Load("env"); err != nil {
		log.Fatalf("Failed to load dotenv. error:%v", err)
	}

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
	router.HandleFunc("/linenotify/on", controller.LinenotifyOn).Methods("GET")
	router.HandleFunc("/linenotify/off", controller.LinenotifyOff).Methods("GET")

	n.UseHandler(router)
	http.Handle("/", n)

	log.Println("init app")
}
