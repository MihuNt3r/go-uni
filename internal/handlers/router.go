package handlers

import (
	"database/sql"
	"net/http"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"

	"go-uni/internal/env"
	"go-uni/internal/repository"
	"go-uni/pkg/middleware"
)

func NewRouter(db *sql.DB) http.Handler {
	studentsRepo := repository.NewStudentsRepository(db)
	teachersRepo := repository.NewTeachersRepository(db)
	coursesRepo := repository.NewCoursesRepository(db)
	enrollmentsRepo := repository.NewEnrollmentsRepository(db)

	studentsHandler := NewStudentsHandler(studentsRepo)
	teachersHandler := NewTeachersHandler(teachersRepo)
	coursesHandler := NewCoursesHandler(coursesRepo)
	enrollmentsHandler := NewEnrollmentsHandler(enrollmentsRepo)
	authHandler := NewAuthHandler(
		env.GetString("AUTH_USERNAME", "admin"),
		env.GetString("AUTH_PASSWORD", "admin"),
		env.GetString("AUTH_JWT_SECRET", "change-me"),
		time.Duration(env.GetInt("AUTH_TOKEN_TTL_MINUTES", 60))*time.Minute,
	)

	authMiddleware := middleware.AuthMiddleware(middleware.AuthConfig{
		JWTSecret: env.GetString("AUTH_JWT_SECRET", "change-me"),
	})

	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/token", authHandler.Token)

	mux.HandleFunc("GET /students", studentsHandler.List)
	mux.HandleFunc("GET /students/{id}", studentsHandler.Get)
	mux.Handle("POST /students", authMiddleware(http.HandlerFunc(studentsHandler.Create)))
	mux.Handle("PUT /students/{id}", authMiddleware(http.HandlerFunc(studentsHandler.Update)))
	mux.Handle("DELETE /students/{id}", authMiddleware(http.HandlerFunc(studentsHandler.Delete)))

	mux.HandleFunc("GET /teachers", teachersHandler.List)
	mux.HandleFunc("GET /teachers/{id}", teachersHandler.Get)
	mux.Handle("POST /teachers", authMiddleware(http.HandlerFunc(teachersHandler.Create)))
	mux.Handle("PUT /teachers/{id}", authMiddleware(http.HandlerFunc(teachersHandler.Update)))
	mux.Handle("DELETE /teachers/{id}", authMiddleware(http.HandlerFunc(teachersHandler.Delete)))

	mux.HandleFunc("GET /courses", coursesHandler.List)
	mux.HandleFunc("GET /courses/{id}", coursesHandler.Get)
	mux.Handle("POST /courses", authMiddleware(http.HandlerFunc(coursesHandler.Create)))
	mux.Handle("PUT /courses/{id}", authMiddleware(http.HandlerFunc(coursesHandler.Update)))
	mux.Handle("DELETE /courses/{id}", authMiddleware(http.HandlerFunc(coursesHandler.Delete)))

	mux.Handle("POST /students/{id}/courses/{course_id}", authMiddleware(http.HandlerFunc(enrollmentsHandler.Enroll)))
	mux.Handle("DELETE /students/{id}/courses/{course_id}", authMiddleware(http.HandlerFunc(enrollmentsHandler.Unenroll)))

	mux.Handle("GET /swagger/", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))
	mux.Handle("GET /swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))

	return middleware.RequestLoggingMiddleware(mux)
}
