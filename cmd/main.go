package main

import (
	"database/sql"
	"go_rest_mysql/config"
	"go_rest_mysql/middleware"
	"go_rest_mysql/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var db *sql.DB

func main() {
	// Initialize the database connection
	db := config.ConnectDB()
	log.Println("Connected to the database")
	defer db.Close() // Close the database connection when the program exits

	// Pass the `db` connection to your controllers as needed
	routes.InitializeControllers(db) // This function will initialize db in controllers (discussed below)

	// Initialize the router
	r := mux.NewRouter()

	// Register routes
	routes.UserRoutes(r)
	// Protected routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTAuthMiddleware)
	routes.CrudRoutes(api) // Apply middleware to protect CRUD routes

	// Start the server
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
