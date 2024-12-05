package main

import (
	"os"
	"time"
	"users/pubsub"
	"log"
	"users/api"
	"users/database"
)

var channels = [...]string{
    "chan1",
    "chan2",
}

func main() {
	db, err := database.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to init db. Error: %s\n", err.Error())
	}

    sendCh := make(chan pubsub.Message, 100)
    receiveCh := make(chan pubsub.Message, 100)

    go pubsub.Transmitter(sendCh)
    for _, ch := range(channels) {
        go pubsub.Listener(receiveCh, ch)
    }
    go pubsub.Interpreter(receiveCh)

    time.Sleep(1 *time.Second)

    if os.Args[1] == "1" {
        sendCh <- pubsub.Message{
            Payload: "This here is a message",
            Channel: channels[0],
        }
    }

	api.InitRouter(db).Run()
}
