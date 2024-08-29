package main

import (
	"blog-aggregator/backend/internal/database"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string       `json:"id"`
	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
	Name      string       `json:"name"`
	ApiKey    string       `json:"api_key"`
	Password  string       `json:"password"`
}

func createUser(cfg *apiConfig, w http.ResponseWriter, r *http.Request) {
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
	b := make([]byte, 64)
	_, err1 := rand.Read(b)
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	h := sha256.New()
	h.Write([]byte(b))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	user.ApiKey = sha1_hash

	// Encode the password
	encPW, err1 := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err1 != nil {
		http.Error(w, "Could not use password", http.StatusInternalServerError)
		return
	}

	// Step 3: Insert into the database
	ctx := r.Context()
	_, err := cfg.DB.CreateUser(ctx, database.CreateUserParams{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Apikey:    user.ApiKey,
		Password:  string(encPW),
	})
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Respond with the created user
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

}

func loginUser(cfg *apiConfig, w http.ResponseWriter, r *http.Request) error {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return err
	}

	ctx := r.Context()
	u, err := cfg.DB.GetUser(ctx, user.Name)
	if err != nil {
		http.Error(w, "Couldn't find user", http.StatusNotFound)
		return err
	}

	// Decrypt found user's password and compare it
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return err
	}

	user.ApiKey = u.Apikey
	user.ID = u.ID
	user.CreatedAt = u.CreatedAt
	user.UpdatedAt = u.UpdatedAt
	user.Name = u.Name
	user.Password = u.Password

	token, err := jwtCreation(user, cfg.jwtSecret)
	if err != nil {
		return err
	}

	response := map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"token": token,
	}

	cfg.jwtToken = token
	setCookieHandler(w, r, cfg, user.ID)
	json.NewEncoder(w).Encode(response)
	return nil
}

func getUserHelper(cfg *apiConfig, w http.ResponseWriter, r *http.Request, userID string) User {
	var user User
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return user
	}
	ctx := r.Context()
	u, err := cfg.DB.GetUserByID(ctx, userID)
	if err != nil {
		http.Error(w, "Issue getting user", http.StatusInternalServerError)
		return user
	}
	user.ApiKey = u.Apikey
	user.ID = u.ID
	user.CreatedAt = u.CreatedAt
	user.UpdatedAt = u.UpdatedAt
	user.Name = u.Name
	user.Password = u.Password

	return user
}
