package routes

import (
	"github.com/example/hrms-backend/controllers"
	"github.com/example/hrms-backend/middlewares"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func registerAttendanceRoutes(r *mux.Router, db *gorm.DB) {
	c := controllers.NewAttendanceController(db)
	s := r.PathPrefix("/attendance").Subrouter()
	s.Use(middlewares.JWTAuth)

	// Employee self-service routes
	s.HandleFunc("", c.ListMine).Methods("GET")
	s.HandleFunc("", c.AddMine).Methods("POST")
	s.HandleFunc("/{id:[0-9]+}", c.DeleteMine).Methods("DELETE")

	// HR
	hr := s.NewRoute().Subrouter()
	hr.Use(middlewares.RequireRole("HR"))
	hr.HandleFunc("", c.ListAll).Methods("GET")
	hr.HandleFunc("/{id:[0-9]+}", c.UpdateAny).Methods("PUT")
}
