package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/example/hrms-backend/docs" // swagger docs generated at build

	"github.com/example/hrms-backend/config"
	"github.com/example/hrms-backend/routes"
)

// @title HRMS Backend API
// @version 1.0
// @description HRMS backend with JWT, RBAC, and concurrency-safe operations.
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Load env
	_ = godotenv.Load()

	// Connect DB
	db, err := config.Connect()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	sqlDB, _ := db.DB()
	if err = sqlDB.Ping(); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}
	log.Println("database connected")

	// Auto migrations
	if err := config.AutoMigrate(db); err != nil {
		log.Fatalf("auto-migration failed: %v", err)
	}

	r := mux.NewRouter()

	// Simple CORS
	r.Use(simpleCORS)

	// Swagger UI using embedded JSON
	swaggerJSON := `{
		"swagger": "2.0",
		"info": {
			"title": "HRMS Backend API",
			"version": "1.0"
		},
		"basePath": "/api/v1",
		"schemes": ["http"],
		"paths": {
			"/auth/register": {
				"post": {
					"summary": "Register",
					"tags": ["Auth"],
					"parameters": [{"in": "body", "name": "body", "required": true, "schema": {"type": "object"}}],
					"responses": {"201": {"description": "created"}}
				}
			},
			"/auth/login": {
				"post": {
					"summary": "Login",
					"tags": ["Auth"],
					"parameters": [{"in": "body", "name": "body", "required": true, "schema": {"type": "object"}}],
					"responses": {"200": {"description": "ok"}}
				}
			},
			"/employees": {
				"get": {
					"summary": "List employees",
					"tags": ["Employees"],
					"security": [{"BearerAuth": []}],
					"responses": {"200": {"description": "ok"}}
				},
				"post": {
					"summary": "Create employee",
					"tags": ["Employees"],
					"security": [{"BearerAuth": []}],
					"parameters": [{"in": "body", "name": "body", "required": true, "schema": {"type": "object"}}],
					"responses": {"201": {"description": "created"}}
				}
			},
			"/attendance": {
				"get": {
					"summary": "List my attendance",
					"tags": ["Attendance"],
					"security": [{"BearerAuth": []}],
					"responses": {"200": {"description": "ok"}}
				},
				"post": {
					"summary": "Add my attendance",
					"tags": ["Attendance"],
					"security": [{"BearerAuth": []}],
					"parameters": [{"in": "body", "name": "body", "required": true, "schema": {"type": "object"}}],
					"responses": {"201": {"description": "created"}}
				}
			},
			"/leaves": {
				"get": {
					"summary": "List my leaves",
					"tags": ["Leaves"],
					"security": [{"BearerAuth": []}],
					"responses": {"200": {"description": "ok"}}
				},
				"post": {
					"summary": "Apply leave",
					"tags": ["Leaves"],
					"security": [{"BearerAuth": []}],
					"parameters": [{"in": "body", "name": "body", "required": true, "schema": {"type": "object"}}],
					"responses": {"201": {"description": "created"}}
				}
			}
		},
		"securityDefinitions": {
			"BearerAuth": {
				"type": "apiKey",
				"name": "Authorization",
				"in": "header"
			}
		}
	}`

	// Swagger routes
	r.Handle("/docs/doc.json", http.FileServer(http.Dir("./docs")))
	r.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/docs/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(swaggerJSON))
	}).Methods("GET")
	r.HandleFunc("/docs/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(swaggerJSON))
	}).Methods("GET")

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()
	routes.Register(api, db)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8082"
	}
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       20 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("server starting on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
	log.Println("server gracefully stopped")
}

func simpleCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
