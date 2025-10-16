package routes

import (
    "github.com/gorilla/mux"
    "gorm.io/gorm"
    "github.com/example/hrms-backend/controllers"
    "github.com/example/hrms-backend/middlewares"
)

func registerEmployeeRoutes(r *mux.Router, db *gorm.DB) {
    c := controllers.NewEmployeeController(db)
    s := r.PathPrefix("/employees").Subrouter()
    s.Use(middlewares.JWTAuth)
    // HR routes
    hr := s.NewRoute().Subrouter()
    hr.Use(middlewares.RequireRole("HR"))
    hr.HandleFunc("", c.List).Methods("GET")
    hr.HandleFunc("", c.Create).Methods("POST")
    hr.HandleFunc("/{id}", c.Update).Methods("PUT")
    hr.HandleFunc("/{id}", c.Delete).Methods("DELETE")
    // Employee self
    s.HandleFunc("/me", c.GetMe).Methods("GET")
}


