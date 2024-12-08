package main

import (
	"log"
	"primary/api"
	"primary/database"
)

func main() {
	db, err := database.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to init db. Error: %s\n", err.Error())
	}

	api.InitRouter(db).Run()
}
