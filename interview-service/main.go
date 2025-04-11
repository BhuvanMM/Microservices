package main

import (
	"context"
	"devlink/common"
	"encoding/json"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	JobID     string `json:"job_id"`
	CoverNote string `json:"cover_note"`
}

type Interview struct {
	ApplicationID string    `json:"application_id"`
	UserID        string    `json:"user_id"`
	JobID         string    `json:"job_id"`
	ScheduledAt   time.Time `json:"scheduled_at"`
}

var collection *mongo.Collection
var rabbitChannel *amqp091.Channel

func init() {
	client := common.GetMongoClient("mongodb://localhost:27017")
	collection = client.Database("devlink_interview").Collection("interviews")
	rabbitChannel = common.ConnectRabbitMQ()

	_, err := rabbitChannel.QueueDeclare("interview_queue", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare interview_queue:", err)
	}

	_, err = rabbitChannel.QueueDeclare("notification_queue", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare notification_queue:", err)
	}
}

func scheduleInterview(app Application) {
	interview := Interview{
		ApplicationID: app.ID,
		UserID:        app.UserID,
		JobID:         app.JobID,
		ScheduledAt:   time.Now().Add(48 * time.Hour), // interview 2 days later
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, interview)

	if err != nil {
		log.Println("Failed to schedule interview:", err)
		return
	}

	body, _ := json.Marshal(interview)
	err = rabbitChannel.Publish(
		"",                   // exchange
		"notification_queue", // queue name
		false, false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Println("Failed to send notification:", err)
	}

}
func consumeApplications() {
	msgs, err := rabbitChannel.Consume(
		"interview_queue", "", true, false, false, false, nil,
	)

	if err != nil {
		log.Fatal("Failed to consume:", err)
	}

	for msg := range msgs {
		var app Application
		if err := json.Unmarshal(msg.Body, &app); err != nil {
			log.Println("Invalid message:", err)
			continue
		}
		scheduleInterview(app)
		log.Println("Interview scheduled for:", app.UserID)
	}
}

func main() {
	log.Println("Interview Service listening for job applications...")
	consumeApplications()
	select {}
}
