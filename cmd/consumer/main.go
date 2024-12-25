package main

import (
	"context"
	"log"
	"time"

	"github.com/dotpep/golang-event-driven-rabbitmq/internal"
	"golang.org/x/sync/errgroup"
)

func main() {
	log.Println("Starting Consumer Logic...")

	log.Println("Connecting to the RabbitMQ...")
	conn, err := internal.ConnectRabbitMQ("admin", "admin", "localhost:5672", "customers")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	log.Println("New Client RabbitMQ - AMQP TCP Channel")
	client, err := internal.NewRabbitMQClient(conn)
	if err != nil {
		panic(err)
	}

	defer client.Close()

	log.Println("Starting Consume... of email-service for Queue - customers-created")
	messageBus, err := client.Consume("customers_created", "email-service", false)
	if err != nil {
		panic(err)
	}

	log.Println("Locking... of Golang Channel of AMQP Delivery")
	var blocking chan struct{}

	// set a timeout for 15 secs
	ctx := context.Background()

	log.Println("Creating Golang Context... with Background of 15-second")
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	// errgroup allows us concurrent tasks
	workerGoroutineCount := 10
	g.SetLimit(workerGoroutineCount)

	go func() {
		for message := range messageBus {
			// spawn a worker
			msg := message
			g.Go(func() error {
				log.Printf("New Message: %v", msg)
				//time.Sleep(10 * time.Second)
				if err := msg.Ack(false); err != nil {
					log.Println("Ack message failed!")
					return err
				}
				log.Printf("Acknowledge message %s\n", msg.MessageId)
				return nil
			})
		}
	}()

	log.Println("Consuming... to close the program press CTRL+C")

	log.Println("Blocking... of Golang Channel of AMQP Delivery")
	<-blocking
}
