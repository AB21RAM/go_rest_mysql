package routes

import (
	"database/sql"
	"go_rest_mysql/controllers"
)

func InitializeControllers(db *sql.DB) {
	controllers.InitializeUserController(db) // Pass db to UserController
	controllers.InitializeCrudController(db) // Pass db to CRUD Controller
}
