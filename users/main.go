package main

import (
	"log"
	"users/api"
	"users/database"
)

func main() {
	db, err := database.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to init db. Error: %s\n", err.Error())
	}
	api.InitRouter(db).Run()
}
