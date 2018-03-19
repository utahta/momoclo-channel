package batch

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/utahta/momoclo-channel/adapter/handler/batch"
	"github.com/utahta/momoclo-channel/adapter/handler/middleware"
	"github.com/utahta/momoclo-channel/config"
)

func init() {
	config.MustLoad("config/deploy.toml")

	router := chi.NewRouter()
	router.Use(middleware.AEContext)

	router.Get("/_ah/start", func(w http.ResponseWriter, req *http.Request) {})
	router.Post("/tweet", batch.Tweet)

	http.Handle("/", router)
}
