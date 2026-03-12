package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	_ "go-uni/docs"
	"go-uni/internal/env"
	"go-uni/internal/handlers"
)

// @title go-uni API
// @version 1.0
// @description University REST API for students, teachers, courses, and enrollments.
// @BasePath /
func main() {
	dbAddr := env.GetString("DB_ADDR", "postgres://admin:postgres@localhost:5432/uni_db?sslmode=disable")

	db, err := sql.Open("postgres", dbAddr)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

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
