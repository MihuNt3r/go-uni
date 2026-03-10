package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	_ "go-uni/docs"
	"go-uni/internal/handlers"
)

// @title go-uni API
// @version 1.0
// @description University REST API for students, teachers, courses, and enrollments.
// @BasePath /
func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if pingErr := db.Ping(); pingErr != nil {
		log.Fatalf("failed to ping database: %v", pingErr)
	}

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

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
