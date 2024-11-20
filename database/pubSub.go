package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func Sender() {
    addr := fmt.Sprintf("%s:6379", os.Getenv("PS_HOST"))
	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: addr, // Redis server address
	})

	// Publish a message
    err := rdb.Publish(ctx, "my-channel", "HALOOOOOOOOO").Err()
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Message sent")
}

func Receiver() {
    addr := fmt.Sprintf("%s:6379", os.Getenv("PS_HOST"))
    log.Printf("Address is %s\n", addr)
	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: addr, // Redis server address
	})

	// Subscribe to a channel
	pubsub := rdb.Subscribe(ctx, "my-channel")
	defer pubsub.Close()

	// Goroutine to handle messages
	go func() {
		for msg := range pubsub.Channel() {
			fmt.Printf("Received message: %s\n", msg.Payload)
		}
	}()

    select {}
}
