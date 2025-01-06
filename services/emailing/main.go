package main

import (
	"emailService/pubsub"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
    if err := godotenv.Load(".env"); err != nil {
        log.Fatalf("Failed to load .env. Err: %s\n", err.Error())
    }
    receiveCh := make(chan pubsub.Message, 100)
    channelName := os.Getenv("EMAIL_CHANNEL")
    go pubsub.Listener(receiveCh, channelName)
    go pubsub.Interpreter(receiveCh)

    select {}
}
