package main

import (
	"log"
	"os"

	"github.com/utahta/momoclo-channel/grpc/line/server"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Faild to load env. error:%v", err)
	}

	if err := server.Run(os.Getenv("SERVER_PORT")); err != nil {
		log.Fatalf("Failed to serv. error:%v", err)
	}
}
