package routes

import (
	"github.com/example/hrms-backend/controllers"
	"github.com/example/hrms-backend/middlewares"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func registerLeaveRoutes(r *mux.Router, db *gorm.DB) {
	c := controllers.NewLeaveController(db)
	s := r.PathPrefix("/leaves").Subrouter()
	s.Use(middlewares.JWTAuth)

	// Employee self-service routes
	s.HandleFunc("", c.Apply).Methods("POST")
	s.HandleFunc("", c.ListMine).Methods("GET")
	s.HandleFunc("/{id:[0-9]+}", c.DeleteMine).Methods("DELETE")

	// HR
	hr := s.NewRoute().Subrouter()
	hr.Use(middlewares.RequireRole("HR"))
	hr.HandleFunc("", c.ListAll).Methods("GET")
	hr.HandleFunc("/{id:[0-9]+}/approve", c.Approve).Methods("POST")
	hr.HandleFunc("/{id:[0-9]+}/reject", c.Reject).Methods("POST")
}
