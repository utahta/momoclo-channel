package app

import (
	"github.com/joho/godotenv"
	"github.com/utahta/momoclo-channel/log"
)

func init() {
	if err := godotenv.Load("env"); err != nil {
		log.Basic().Fatalf("Failed to load dotenv. error:%v", err)
	}

	initRoutes()

	log.Basic().Infof("init app")
}
