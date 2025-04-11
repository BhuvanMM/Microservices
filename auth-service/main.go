package main

import (
	"devlink/common"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var users = map[string]User{}

func register(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	user.ID = uuid.NewString()
	users[user.Username] = user
	w.Write([]byte("User registered successfully!"))
}

func login(w http.ResponseWriter, r *http.Request) {
	var creds User
	json.NewDecoder(r.Body).Decode(&creds)
	user, ok := users[creds.Username]
	if !ok || user.Password != creds.Password {
		http.Error(w, "Invalid User!", http.StatusUnauthorized)
		return
	}

	token, _ := common.GenerateJWT(user.ID)
	w.Write([]byte(token))
}

func main() {
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.ListenAndServe(":8001", nil)
}
