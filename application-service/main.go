package main

import (
	"context"
	"devlink/common"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Application struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	JobID     string `json:"job_id"`
	CoverNote string `json:"cover_note"`
}

var collection *mongo.Collection
var rabbitChannel *amqp091.Channel

func init() {
	client := common.GetMongoClient("mongodb://localhost:27017")
	collection = client.Database("devlink_application").Collection("applications")
	rabbitChannel = common.ConnectRabbitMQ()

	_, err := rabbitChannel.QueueDeclare("interview_queue", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare queue:", err)
	}
}

func applyJob(w http.ResponseWriter, r *http.Request) {
	var app Application
	json.NewDecoder(r.Body).Decode(&app)
	app.ID = uuid.NewString()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, app)
	if err != nil {
		http.Error(w, "Failed to apply", http.StatusInternalServerError)
		return
	}

	body, _ := json.Marshal(app)
	err = rabbitChannel.Publish(
		"",
		"interview_queue",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		http.Error(w, "Failed to queue interview", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(app)

}

func list(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"_id": userID})
	if err != nil {
		http.Error(w, "Failed to get applications", http.StatusInternalServerError)
		return
	}

	var apps []Application
	cursor.All(ctx, &apps)
	json.NewEncoder(w).Encode(apps)
}

func main() {
	http.HandleFunc("/apply", applyJob)
	http.HandleFunc("/list", list)
	http.ListenAndServe(":8004", nil)
}
