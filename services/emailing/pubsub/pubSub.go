package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"
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
        m := Message{}
        err := json.Unmarshal([]byte(msg.Payload), &m)
        if err != nil {
            log.Println("Failed to unmarshal msg")
            continue
        }
        receiveCh <- Message{
            Payload: m.Payload,
            Channel: m.Channel,
            Subject: m.Subject,
            Receiver: m.Receiver,
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
    var err error
    for {
        msg = <- sendCh
        err = rdb.Publish(ctx, msg.Channel, msg.Payload).Err()
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
        go sendEmail(msg)
    }
}

func sendEmail(msg Message) {
    mail := gomail.NewMessage()
    mail.SetHeader("From", os.Getenv("EMAIL_ADDRESS"))
    mail.SetHeader("To", msg.Receiver)
    mail.SetHeader("Subject", msg.Subject)
    mail.SetBody("text/plain", msg.Payload)
    log.Printf("Msg: %+v\n", msg)

    dialer := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("EMAIL_ADDRESS"), os.Getenv("EMAIL_PASSWORD"))

    err := dialer.DialAndSend(mail)
    if err != nil {
        log.Printf("Err %s\n", err.Error())
        return
    }
    log.Println("Message sent")
}
