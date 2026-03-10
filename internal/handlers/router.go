package handlers

import (
	"database/sql"
	"net/http"

	"go-uni/internal/repository"
)

func NewRouter(db *sql.DB) *http.ServeMux {
	studentsRepo := repository.NewStudentsRepository(db)
	teachersRepo := repository.NewTeachersRepository(db)
	coursesRepo := repository.NewCoursesRepository(db)
	enrollmentsRepo := repository.NewEnrollmentsRepository(db)

	studentsHandler := NewStudentsHandler(studentsRepo)
	teachersHandler := NewTeachersHandler(teachersRepo)
	coursesHandler := NewCoursesHandler(coursesRepo)
	enrollmentsHandler := NewEnrollmentsHandler(enrollmentsRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /students", studentsHandler.List)
	mux.HandleFunc("GET /students/{id}", studentsHandler.Get)
	mux.HandleFunc("POST /students", studentsHandler.Create)
	mux.HandleFunc("PUT /students/{id}", studentsHandler.Update)
	mux.HandleFunc("DELETE /students/{id}", studentsHandler.Delete)

	mux.HandleFunc("GET /teachers", teachersHandler.List)
	mux.HandleFunc("GET /teachers/{id}", teachersHandler.Get)
	mux.HandleFunc("POST /teachers", teachersHandler.Create)
	mux.HandleFunc("PUT /teachers/{id}", teachersHandler.Update)
	mux.HandleFunc("DELETE /teachers/{id}", teachersHandler.Delete)

	mux.HandleFunc("GET /courses", coursesHandler.List)
	mux.HandleFunc("GET /courses/{id}", coursesHandler.Get)
	mux.HandleFunc("POST /courses", coursesHandler.Create)
	mux.HandleFunc("PUT /courses/{id}", coursesHandler.Update)
	mux.HandleFunc("DELETE /courses/{id}", coursesHandler.Delete)

	mux.HandleFunc("POST /students/{id}/courses/{course_id}", enrollmentsHandler.Enroll)
	mux.HandleFunc("DELETE /students/{id}/courses/{course_id}", enrollmentsHandler.Unenroll)

	return mux
}
