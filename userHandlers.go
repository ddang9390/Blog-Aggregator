package main

import (
	"blog-aggregator/internal/database"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string       `json:"id"`
	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
	Name      string       `json:"name"`
	ApiKey    string       `json:"api_key"`
}

func createUser(cfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		// // Step 1: Parse the request body
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Step 2: Fill additional fields
		user.ID = uuid.New().String()
		user.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
		user.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

		// Make API key
		randInt, _ := rand.Read(make([]byte, 64))
		h := sha256.New()
		h.Write([]byte(strconv.Itoa(randInt)))
		sha1_hash := hex.EncodeToString(h.Sum(nil))
		user.ApiKey = sha1_hash

		// Step 3: Insert into the database
		ctx := r.Context()
		_, err := cfg.DB.CreateUser(ctx, database.CreateUserParams{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Name:      user.Name,
		})
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		// // Respond with the created user
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

func getUser(cfg *apiConfig, w http.ResponseWriter, r *http.Request) (User, error) {
	var user User
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authroization header required", http.StatusUnauthorized)
		return user, fmt.Errorf("uthroization header required")
	}

	apiString := r.Header.Get("ApiKey")
	if apiString == "" {
		http.Error(w, "Api key required", http.StatusUnauthorized)
		return user, fmt.Errorf("uthroization header required")
	}

	fmt.Println(apiString)
	ctx := r.Context()
	u, err := cfg.DB.GetUser(ctx, apiString)
	if err != nil {
		http.Error(w, "Couldn't find user", http.StatusNotFound)
		return user, err
	}
	user.ApiKey = u.Apikey
	user.ID = u.ID
	user.CreatedAt = u.CreatedAt
	user.UpdatedAt = u.UpdatedAt
	user.Name = u.Name

	json.NewEncoder(w).Encode(u)
	return user, nil
}

func getUserHelper(cfg *apiConfig, w http.ResponseWriter, r *http.Request) *User {
	apiString := r.Header.Get("ApiKey")
	if apiString == "" {
		http.Error(w, "Api key required", http.StatusUnauthorized)
		return nil
	}
	user, err := getUser(cfg, w, r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Issue getting user", http.StatusInternalServerError)
		return nil
	}
	return &user
}
