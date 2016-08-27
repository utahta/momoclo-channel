package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/utahta/momoclo-channel/line/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Faild to load env. error:%v", err)
	}

	var (
		channelID     int64
		channelSecret = os.Getenv("LINEBOT_CHANNEL_SECRET")
		channelMID    = os.Getenv("LINEBOT_CHANNEL_MID")
	)
	channelID, err := strconv.ParseInt(os.Getenv("LINEBOT_CHANNEL_ID"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	s, err := server.New(channelID, channelSecret, channelMID)
	if err != nil {
		log.Fatalf("Failed to init line server. error:%v", err)
	}

	if err := s.Run(os.Getenv("SERVER_PORT")); err != nil {
		log.Fatalf("Failed to serv. error:%v", err)
	}
}
