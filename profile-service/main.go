package main

import (
	"context"
	"devlink/common"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Profile struct {
	ID       string   `json:"id"`
	UserID   string   `json:"user_id"`
	Bio      string   `json:"bio"`
	Skills   []string `json:"skills"`
	Projects []string `json:"projects"`
}

var collection *mongo.Collection

func init() {
	client := common.GetMongoClient("mongodb://localhost:27017")
	collection = client.Database("devlink_profile").Collection("profiles")
}

func createProfile(w http.ResponseWriter, r *http.Request) {
	var profile Profile
	json.NewDecoder(r.Body).Decode(&profile)
	profile.ID = uuid.NewString()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, profile)
	if err != nil {
		http.Error(w, "Failed to create profile!", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(profile)
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userID := r.URL.Query().Get("user_id")

	var profile Profile
	err := collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&profile)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(profile)
}

func updateProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	var profile Profile
	json.NewDecoder(r.Body).Decode(&profile)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{
		"$set": bson.M{
			"bio":     profile.Bio,
			"skills":  profile.Skills,
			"project": profile.Projects,
		},
	})

	if err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Profile updated"))
}

func deleteProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"user_id": userID})

	if err != nil {
		http.Error(w, "Failed to delete", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Profile deleted"))
}

func main() {
	http.HandleFunc("/create", createProfile)
	http.HandleFunc("/get", getProfile)
	http.HandleFunc("/update", updateProfile)
	http.HandleFunc("/delete", deleteProfile)
	http.ListenAndServe(":8002", nil)
}
