package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

type Producer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
}

func NewProducer(rabbitMQURL, queueName string) (*Producer, error) {
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

	return &Producer{
		connection: conn,
		channel:    ch,
		queueName:  queueName,
	}, nil
}

func (p *Producer) Publish(message []byte) error {
	err := p.channel.Publish(
		"",
		p.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		log.Printf("Error publishing message to RabbitMQ: %v", err)
	}
	return err
}

func (p *Producer) Close() {
	p.channel.Close()
	p.connection.Close()
}
