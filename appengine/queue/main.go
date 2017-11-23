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
	router.Use(middleware.AppengineContext)

	router.Get("/_ah/start", func(w http.ResponseWriter, req *http.Request) {})
	router.Post("/queue/tweet", queue.Tweet)
	router.Post("/queue/line/broadcast", queue.LineNotifyBroadcast)
	router.Post("/queue/line", queue.LineNotify)

	http.Handle("/", router)
}
