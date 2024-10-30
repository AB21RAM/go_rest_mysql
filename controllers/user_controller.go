package controllers

import (
	"database/sql"
	"encoding/json"
	"go_rest_mysql/models"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func InitializeUserController(database *sql.DB) {
	db = database
}

var jwtSecret = []byte("98991a0ce985e1e48f796693b95a33e52511cbc993e24ccbd0fdbc91cec73e21b37cfe5f50b35c002346fc06406104745b3c5f1426d81c632311be1a72e00743d044cb37af3024af41fb746a4daf908a5ba12c3f34805f18a4e8229026e70278916542bfe0475e6c7fd765c928988b29ce1d27876add3b284f04ef330ddb8b266fffa34191790000d0bc19a4ac1c28276363173a046e60015d7777bd3f04b682995f0a4f7c49eead65e1819a4e548108ae32b280228c3ec95827b8f0d9f645717afbe2d85412ebf1aa9cba7920c7a39c78bfe5f512ab493f006eb7205ab315cd509a6849fed1addb7c76bff2b821b17d91c5db2b7185f9f5644dfd66a9a5e0aa")

func GenerateJWT(email int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})
	return token.SignedString(jwtSecret)
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
		// user already present with the same mail
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]bool{"exits": true})
		// http.Error(w, "Failed to create user", http.StatusInternalServerError)
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
	// Generate JWT token
	token, err := GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
	// // On successful login, send a success response
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(map[string]string{"message": "SUCCESS"})
}
