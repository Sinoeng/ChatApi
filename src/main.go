package main

import (
	"log"
	"primary/api"
	"primary/database"
	"primary/pubsub"
	"primary/utils"
)

func main() {
	db, err := database.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to init db. Error: %s\n", err.Error())
	}
    utils.CreateDefaultAdmin(db)

    sendCh := make(chan pubsub.Message, 100)
    go pubsub.Transmitter(sendCh)
	api.InitRouter(db, sendCh).Run()
}
