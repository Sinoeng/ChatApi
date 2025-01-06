package main

import (
	"emailService/pubsub"
	"os"
)

func main() {
    receiveCh := make(chan pubsub.Message, 100)
    channelName := os.Getenv("EMAIL_CHANNEL")
    go pubsub.Listener(receiveCh, channelName)
    go pubsub.Interpreter(receiveCh)

    select {}
}
