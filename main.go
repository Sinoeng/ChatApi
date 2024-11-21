package main

import (
	"ChatApi/database"
	"log"
	"os"
	"time"
)


func main() {
    para := os.Args[1]
    log.Printf("Parameter: %s\n", para)
    switch para {
    case "1":
        log.Println("Starting sender")
        for {
            database.Sender()
            time.Sleep(1 * time.Second)
        }
    case "2":
        log.Println("Starting receiver")
        database.Receiver(os.Args[2])
    default:
        log.Fatalln("In the default case")
    }
}
