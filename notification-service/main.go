package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"devlink/common"

	"github.com/rabbitmq/amqp091-go"
)

type Interview struct {
	ApplicationID string    `json:"application_id"`
	UserID        string    `json:"user_id"`
	JobID         string    `json:"job_id"`
	ScheduledAt   time.Time `json:"scheduled_at"`
}

var rabbitChannel *amqp091.Channel

func init() {
	rabbitChannel = common.ConnectRabbitMQ()

	// Declare queue
	_, err := rabbitChannel.QueueDeclare("notification_queue", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare notification_queue:", err)
	}
}

func listenForNotifications() {
	msgs, err := rabbitChannel.Consume(
		"notification_queue", "", true, false, false, false, nil,
	)

	if err != nil {
		log.Fatal("Failed to consume notification messages:", err)
	}

	for msg := range msgs {
		var interview Interview
		if err := json.Unmarshal(msg.Body, &interview); err != nil {
			log.Println("Invalid message:", err)
			continue
		}
		fmt.Printf("📢 [NOTIFY] User %s: Your interview for Job %s is scheduled on %s\n",
			interview.UserID, interview.JobID, interview.ScheduledAt.Format("2006-01-02 15:04"))
	}
}

func main() {
	log.Println("Notification Service running...")
	listenForNotifications()
	select {} // block forever
}
