package controllers

import (
	"database/sql"
	"encoding/json"
	"go_rest_mysql/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func InitializeUserController(database *sql.DB) {
	db = database
}

// Register a new user
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Insert the new user into the database
	_, err = db.Exec("INSERT INTO users (name, phone, email, password) VALUES (?, ?, ?, ?)", user.Name, user.Phone, user.Email, user.Password)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "SUCCESS"})
}

// Login user
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Retrieve user from database
	var user models.User
	err = db.QueryRow("SELECT id, name, phone, email, password FROM users WHERE email = ?", creds.Email).Scan(&user.ID, &user.Name, &user.Phone, &user.Email, &user.Password)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Compare the stored hashed password with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// On successful login, send a success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "SUCCESS"})
}
