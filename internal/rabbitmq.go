package internal

import (
	"context"
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

// CreateQueue will create a new Queue based on given cfgs
func (rc RabbitClient) CreateQueue(queueName string, durable, autoDelete bool) error {
	_, err := rc.ch.QueueDeclare(
		queueName, durable, autoDelete,
		false, false, nil,
	)
	return err
}

// CreateBinding will bind the current Channel to the given Exchange using the Routingkey provided
func (rc RabbitClient) CreateBinding(name, binding, exchange string) error {
	// leaving noWait false, having noWait set to false will make the channel return an error if its fales to bind
	return rc.ch.QueueBind(name, binding, exchange, false, nil)
}

// Send is used to Publish payloads onto an Exchange with the given Routingkey
func (rc RabbitClient) Send(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error {
	return rc.ch.PublishWithContext(
		ctx,
		exchange, routingKey,
		// Mandatory is used to determine if an error should be returned upon failure of sending
		true,
		// Immediate is removed in RabbitMQ:3
		false,
		// msg amqp.Publishing - Options is actuall message that we're sending
		options,
	)
}
