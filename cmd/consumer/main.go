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

			// Acknowledge Exchange that Message was Successfully Consumed
			//if err := message.Ack(false); err != nil {
			//	log.Printf("Acknowledge message Failed! %v\n", err)
			//	continue
			//}

			// Nacking RabbitMQ that it was actuall Failure with consuming message
			if !message.Redelivered {
				message.Nack(false, true)
				continue
			}

			if err := message.Ack(false); err != nil {
				log.Println("Failed to ack message")
				continue
			}

			log.Printf("Acknowledge message %s\n", message.MessageId)
		}
	}()

	log.Println("Consuming... to close the program press CTRL+C")

	log.Println("Blocking... of Golang Channel of AMQP Delivery")
	<-blocking
}
