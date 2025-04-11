package common

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func ConnectRabbitMQ() *amqp091.Channel {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Rabbit MQ Connnection failure!")
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Rabbit MQ Channel failure!")
	}
	return ch
}
