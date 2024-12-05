package users

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func Transmitter(sendCh chan string) {
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

}

func Listener(commReceive chan string, channelName string) {
    addr := fmt.Sprintf("%s:6379", os.Getenv("PS_HOST"))
    log.Printf("Address is %s\n", addr)
	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: addr, // Redis server address
	})

	// Subscribe to a channel
	pubsub := rdb.Subscribe(ctx, channelName)
	defer pubsub.Close()

    for msg := range pubsub.Channel() {
        fmt.Printf("Received message: %s\n", msg.Payload)
        commReceive <- msg.Payload
    }

    panic("Listener is exiting\n\n\n")
}

func Interpreter(receiver chan string) {
    for {
        msg := <- receiver
        log.Printf("Interpreter says: %s\n", msg)
    }
}
