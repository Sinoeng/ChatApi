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
    msg := fmt.Sprintf("MESSAGE")
    err := rdb.Publish(ctx, "my-channel", msg).Err()
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Message sent")

    // Publish a message
    msg = fmt.Sprintf("MESSAGE2")
    err = rdb.Publish(ctx, "channel2", msg).Err()
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Message2 sent")
}

func Receiver(arg string) {
    addr := fmt.Sprintf("%s:6379", os.Getenv("PS_HOST"))
    log.Printf("Address is %s\n", addr)
	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: addr, // Redis server address
	})

	// Subscribe to a channel
	pubsub := rdb.Subscribe(ctx, "my-channel")
	defer pubsub.Close()
    pubsub2 := rdb.Subscribe(ctx, "channel2")
    defer pubsub.Close()

	// Goroutine to handle messages
    switch arg {
    case "1":
        go func() {
            for msg := range pubsub.Channel() {
                fmt.Printf("Received message: %s\n", msg.Payload)
            }
        }()
    case "2":
        go func() {
            for msg := range pubsub2.Channel() {
                fmt.Printf("Received message: %s\n", msg.Payload)
            }
        }()
    }

    select {}
}
