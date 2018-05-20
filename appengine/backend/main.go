package backend

import (
	"github.com/utahta/momoclo-channel/api"
	"github.com/utahta/momoclo-channel/config"
)

func init() {
	config.MustLoad("config/deploy.toml")

	s := api.NewBackendServer()
	s.Handle()
}
