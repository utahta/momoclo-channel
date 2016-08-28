package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/utahta/momoclo-channel/linebot/server"
	mlog "github.com/utahta/momoclo-channel/log"
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

	fp, err := os.OpenFile(os.Getenv("LOG_PATH"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open file. error:%v", err)
	}
	defer fp.Close()

	s.Log = mlog.NewIOLogger(fp)

	if err := s.Run(os.Getenv("SERVER_PORT")); err != nil {
		log.Fatalf("Failed to serv. error:%v", err)
	}
}
