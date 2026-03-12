package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	_ "go-uni/docs"
	"go-uni/internal/database"
	"go-uni/internal/env"
	"go-uni/internal/handlers"
)

// @title go-uni API
// @version 1.0
// @description University REST API for students, teachers, courses, and enrollments.
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Use format: Bearer {token}
func main() {
	dbAddr := env.GetString("DB_ADDR", "postgres://admin:postgres@localhost:5432/uni_db?sslmode=disable")

	db, err := database.New(dbAddr, 5, 5)
	if err != nil {
		log.Fatalf("error connecting to database: %s", err)
	}

	addr := env.GetString("HTTP_ADDR", ":8080")

	server := &http.Server{
		Addr:              addr,
		Handler:           handlers.NewRouter(db),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("server listening on %s", addr)
	if serveErr := server.ListenAndServe(); serveErr != nil && serveErr != http.ErrServerClosed {
		log.Fatalf("server failed: %v", serveErr)
	}
}
