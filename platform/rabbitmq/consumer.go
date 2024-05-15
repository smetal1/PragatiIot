package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

type MessageHandler interface {
	HandleMessage(message []byte) error
}

type Consumer struct {
	connection     *amqp.Connection
	channel        *amqp.Channel
	queueName      string
	messageHandler MessageHandler
}

func NewConsumer(rabbitMQURL, queueName string, messageHandler MessageHandler) (*Consumer, error) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,  // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		connection:     conn,
		channel:        ch,
		queueName:      queueName,
		messageHandler: messageHandler,
	}, nil
}

func (c *Consumer) StartConsuming() {
	msgs, err := c.channel.Consume(
		c.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error starting consumer: %v", err)
	}

	go func() {
		for d := range msgs {
			err := c.messageHandler.HandleMessage(d.Body)
			if err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}
	}()
}

func (c *Consumer) Close() {
	c.channel.Close()
	c.connection.Close()
}
