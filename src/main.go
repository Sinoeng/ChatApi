package main

import (
	"log"
	"os"
	"primary/api"
	"primary/database"
	"primary/pubsub"
	"primary/utils"
)

func main() {
	db, err := database.InitDatabase(os.Getenv("MYSQL_DATABASE"))
	if err != nil {
		log.Fatalf("Failed to init db. Error: %s\n", err.Error())
	}
    utils.CreateDefaultAdmin(db)

    sendCh := make(chan pubsub.Message, 100)
    go pubsub.Transmitter(sendCh)
	api.InitRouter(db, sendCh).Run()
}
