package routes

import (
	"go_rest_mysql/controllers"

	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router) {
	r.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
	r.HandleFunc("/login", controllers.LoginUser).Methods("POST")
}
