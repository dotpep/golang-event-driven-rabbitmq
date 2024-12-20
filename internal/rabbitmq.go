package internal

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	// The connection used by the client
	// a connection is TCP connection,
	// you should reuse the connection across your whole application,
	// also spawn new channels on every concurrent task that is running (TCP connection).
	conn *amqp.Connection
	// a channel is a multiplexed connection over the TCP connection,
	// is like separate connection but using the TCP that we setup on connection (Sub connection of TCP).
	// Channel is usedd to process / Send messages.
	ch *amqp.Channel
}

func ConnectRabbitMQ(username, password, host, vhost string) (*amqp.Connection, error) {
	// TODO: add port also!
	return amqp.Dial(
		fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost),
	)
}

func NewRabbitMQClient(conn *amqp.Connection) (RabbitClient, error) {
	// take the connection and spawn a channel from it,
	// and this channel will be use for this created rabbit client,
	// this allows us to reuse the connection between multiple RabbitMQ clients

	ch, err := conn.Channel()
	if err != nil {
		return RabbitClient{}, err
	}

	return RabbitClient{
		conn: conn,
		ch:   ch,
	}, nil
}

// Closes the channel, not connection
func (rc RabbitClient) Close() error {
	return rc.ch.Close()
}

func (rc RabbitClient) CreateQueue(queueName string, durable, autoDelete bool) error {
	_, err := rc.ch.QueueDeclare(
		queueName, durable, autoDelete,
		false, false, nil,
	)
	return err
}
