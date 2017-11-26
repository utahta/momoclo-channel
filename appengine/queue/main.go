package queue

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/utahta/momoclo-channel/adapter/handler/middleware"
	"github.com/utahta/momoclo-channel/adapter/handler/queue"
	"github.com/utahta/momoclo-channel/lib/config"
)

func init() {
	config.MustLoad("config/deploy.toml")

	router := chi.NewRouter()
	router.Use(middleware.AEContext)

	router.Get("/_ah/start", func(w http.ResponseWriter, req *http.Request) {})

	router.Route("/queue", func(r chi.Router) {
		r.Post("/tweet", queue.Tweet)
		r.Post("/line/broadcast", queue.LineNotifyBroadcast)
		r.Post("/line", queue.LineNotify)
	})

	http.Handle("/", router)
}
