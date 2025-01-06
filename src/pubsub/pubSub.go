package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

type Message struct {
    Payload string
    Receiver string
    Subject string
    Channel string
}

var ctx = context.Background()

func Listener(receiveCh chan Message, channelName string) {
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
        receiveCh <- Message{
            Payload: msg.Payload,
            Channel: msg.Channel,
        }
    }

    panic("Listener is exiting\n\n\n")
}

func Transmitter(sendCh chan Message) {
    addr := fmt.Sprintf("%s:6379", os.Getenv("PS_HOST"))
	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: addr, // Redis server address
	})

    // Publish a message
    var msg Message
    for {
        msg = <- sendCh
        m, err := json.Marshal(msg)
        if err != nil {
            log.Println("Error marshalling msg")
            continue
        }
        err = rdb.Publish(ctx, msg.Channel, m).Err()
        if err != nil {
            log.Fatal(err)
        }
        log.Println("Message sent")
    }
}

func Interpreter(receiveCh chan Message) {
    var msg Message
    for {
        msg = <- receiveCh
        log.Printf("Interpreter says: %s\n", msg.Payload)
    }
}
