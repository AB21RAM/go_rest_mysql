package controllers

import (
	"database/sql"
	"encoding/json"
	"go_rest_mysql/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var db *sql.DB

type Info struct {
	Name  string
	Phone string
	Email string
}

func InitializeCrudController(database *sql.DB) {
	db = database
}

// Get all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("ALL USERS CALLED")
	rows, err := db.Query("SELECT name, email , phone FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []Info
	for rows.Next() {
		var user Info
		if err := rows.Scan(&user.Name, &user.Email, &user.Phone); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Get a user by ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	// type Error struct {
	// 	ErrorString string `json:"error"`
	// }
	// var error Error
	log.Println("USER WITH ID: " + id + " CALLED")
	var user Info
	err := db.QueryRow("SELECT name, email, phone FROM users WHERE id = ?", id).Scan(&user.Name, &user.Email, &user.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Not Needed as we are creating the Register feature
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	user.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Update a user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	type UpdateUserBody struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Email string `json:"email"`
	}
	params := mux.Vars(r)
	id := params["id"]

	var user UpdateUserBody
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE users SET name = ?, email = ? , phone = ? WHERE id = ?", user.Name, user.Email, user.Phone, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.ID, _ = strconv.Atoi(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Delete a user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
