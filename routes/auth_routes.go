package routes

import (
    "github.com/gorilla/mux"
    "gorm.io/gorm"
    "github.com/example/hrms-backend/controllers"
)

func registerAuthRoutes(r *mux.Router, db *gorm.DB) {
    c := controllers.NewAuthController(db)
    s := r.PathPrefix("/auth").Subrouter()
    s.HandleFunc("/register", c.Register).Methods("POST")
    s.HandleFunc("/login", c.Login).Methods("POST")
    s.HandleFunc("/refresh", c.Refresh).Methods("POST")
}


