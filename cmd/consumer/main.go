package main

import (
	"log"

	"github.com/dotpep/golang-event-driven-rabbitmq/internal"
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

	go func() {
		for message := range messageBus {
			log.Printf("New Message: %v\n", message)
		}
	}()

	log.Println("Consuming... to close the program press CTRL+C")

	log.Println("Blocking... of Golang Channel of AMQP Delivery")
	<-blocking
}
