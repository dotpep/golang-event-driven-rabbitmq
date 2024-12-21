package main

import (
	"context"
	"log"
	"time"

	"github.com/dotpep/golang-event-driven-rabbitmq/internal"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	log.Println("Start Producer Logic...")

	// TODO: add credentials to .env and load it
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

	log.Println("Creating new Queue... named customers_created")
	if err := client.CreateQueue("customers_created", true, false); err != nil {
		panic(err)
	}

	log.Println("Creating new Queue... named customers_test")
	if err := client.CreateQueue("customers_test", false, true); err != nil {
		panic(err)
	}

	log.Println("Bounding Bind... for an Exchange - customer_events to the Queue - customers_created")
	if err := client.CreateBinding("customers_created", "customers.created.*", "customer_events"); err != nil {
		panic(err)
	}

	log.Println("Bounding Bind... for an Exchange - customer_events to the Queue - customers_test")
	if err := client.CreateBinding("customers_test", "customers.*", "customer_events"); err != nil {
		panic(err)
	}

	log.Println("Creating Golang Context... with Background")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Sending Persistent Message... to the Exchange - customer_events with RoutingKey - customers.created.us")
	if err := client.Send(ctx, "customer_events", "customers.created.us", amqp091.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp091.Persistent,
		Body:         []byte(`An cool message between services`),
	}); err != nil {
		panic(err)
	}

	log.Println("Sending Transient Message... to the Exchange - customer_events with RoutingKey - customers.test")
	if err := client.Send(ctx, "customer_events", "customers.test", amqp091.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp091.Transient,
		Body:         []byte(`An uncool undurable message`),
	}); err != nil {
		panic(err)
	}

	time.Sleep(10 * time.Second)

	log.Println(client)
	defer log.Println("Shutdown a Producer...")
}
