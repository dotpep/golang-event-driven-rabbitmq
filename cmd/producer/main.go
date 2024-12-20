package main

import (
	"log"
	"time"

	"github.com/dotpep/golang-event-driven-rabbitmq/internal"
)

func main() {
	log.Println("Start Producer Logic...")

	// TODO: add credentials to .env and load it
	log.Println("RabbitMQ Connection")
	conn, err := internal.ConnectRabbitMQ("admin", "admin", "localhost:5672", "customers")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	log.Println("RabbitMQ New Client")
	client, err := internal.NewRabbitMQClient(conn)
	if err != nil {
		panic(err)
	}

	defer client.Close()

	log.Println("Queue - customers_created")
	if err := client.CreateQueue("customers_created", true, false); err != nil {
		panic(err)
	}

	log.Println("Queue - customers_test")
	if err := client.CreateQueue("customers_test", false, true); err != nil {
		panic(err)
	}

	time.Sleep(10 * time.Second)

	log.Println(client)
	defer log.Println("Shutdown a Producer...")
}
