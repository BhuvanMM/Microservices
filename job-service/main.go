package main

import (
	"context"
	"devlink/common"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Job struct {
	ID          string `json:"id"`
	Company     string `json:"company"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
}

var collection *mongo.Collection

func init() {
	client := common.GetMongoClient("mongodib://localhost:27017")
	collection = client.Database("devlink_job").Collection("jobs")
}

func createJob(w http.ResponseWriter, r *http.Request) {
	var job Job
	json.NewDecoder(r.Body).Decode(&job)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, job)
	if err != nil {
		http.Error(w, "Failed to create job!", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(job)
}

func getJobs(w http.ResponseWriter, r *http.Request) {
	var jobs []Job

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch jobs", http.StatusInternalServerError)
		return
	}

	if err := cursor.All(ctx, &jobs); err != nil {
		http.Error(w, "Error decoding jobs", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(jobs)
}

func updateJob(w http.ResponseWriter, r *http.Request) {
	var job Job
	json.NewDecoder(r.Body).Decode(&job)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, bson.M{"id": job.ID}, bson.M{
		"$set": bson.M{
			"title":       job.Title,
			"description": job.Description,
			"location":    job.Location,
		},
	})

	if err != nil {
		http.Error(w, "Failed to update job", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Job updated"))
}

func deleteJob(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		http.Error(w, "Failed to delete job", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Job deleted"))
}

func main() {
	http.HandleFunc("/create", createJob)
	http.HandleFunc("/get", getJobs)
	http.HandleFunc("/update", updateJob)
	http.HandleFunc("/delete", deleteJob)
	http.ListenAndServe(":8003", nil)
}
