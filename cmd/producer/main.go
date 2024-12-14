package main

import (
	"log"
	"time"

	"github.com/dotpep/golang-event-driven-rabbitmq/internal"
)

func main() {
	// TODO: add credentials to .env and load it
	conn, err := internal.ConnectRabbitMQ("admin", "admin", "localhost:5672", "customers")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	client, err := internal.NewRabbitMQClient(conn)
	if err != nil {
		panic(err)
	}

	defer client.Close()

	time.Sleep(10 * time.Second)

	log.Println(client)
}
