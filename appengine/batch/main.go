package batch

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/utahta/momoclo-channel/config"
	"github.com/utahta/momoclo-channel/handler"
	"github.com/utahta/momoclo-channel/handler/middleware"
)

func init() {
	config.MustLoad("config/deploy.toml")

	router := chi.NewRouter()
	router.Use(middleware.AEContext)

	router.Get("/_ah/start", func(w http.ResponseWriter, req *http.Request) {})
	router.Post("/tweet", handler.Tweet)

	http.Handle("/", router)
}
