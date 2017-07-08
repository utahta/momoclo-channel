package backend

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"github.com/utahta/momoclo-channel/controller"
	"github.com/utahta/momoclo-channel/lib/config"
	"github.com/utahta/momoclo-channel/middleware"
)

func init() {
	config.MustLoad("config/deploy.toml")

	n := negroni.New()
	n.Use(negroni.HandlerFunc(middleware.Appengine))

	router := mux.NewRouter()
	router.HandleFunc("/queue/tweet", controller.QueueTweet).Methods("POST")

	n.UseHandler(router)
	http.Handle("/", n)

	log.Println("init queue app")
}
