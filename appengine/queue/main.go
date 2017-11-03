package queue

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

	router.Get("/_ah/start", func(w http.ResponseWriter, req *http.Request) {})
	router.Post("/queue/line", controller.QueueLine)
	router.Post("/queue/tweet", controller.QueueTweet)

	http.Handle("/", router)

	log.Println("init queue app")
}
