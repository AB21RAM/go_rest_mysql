package routes

import (
	"database/sql"
	"go_rest_mysql/controllers"

	"github.com/gorilla/mux"
)

var db *sql.DB

func CrudRoutes(r *mux.Router) {
	r.HandleFunc("/users", controllers.GetUsers).Methods("GET")
	r.HandleFunc("/user", controllers.GetUser).Methods("GET")
	r.HandleFunc("/user", controllers.CreateUser).Methods("POST")
	r.HandleFunc("/user/{id}", controllers.UpdateUser).Methods("PUT")
	r.HandleFunc("/user/{id}", controllers.DeleteUser).Methods("DELETE")
}
