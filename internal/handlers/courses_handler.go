package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"go-uni/internal/models"
	"go-uni/internal/repository"
)

type CoursesHandler struct {
	repo *repository.CoursesRepository
}

func NewCoursesHandler(repo *repository.CoursesRepository) *CoursesHandler {
	return &CoursesHandler{repo: repo}
}

func (h *CoursesHandler) List(w http.ResponseWriter, r *http.Request) {
	courses, err := h.repo.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list courses")
		return
	}

	writeJSON(w, http.StatusOK, courses)
}

func (h *CoursesHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid course id")
		return
	}

	course, getErr := h.repo.GetByID(r.Context(), id)
	if getErr != nil {
		writeError(w, http.StatusInternalServerError, "failed to get course")
		return
	}
	if course == nil {
		writeError(w, http.StatusNotFound, "course not found")
		return
	}

	writeJSON(w, http.StatusOK, course)
}

func (h *CoursesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var course models.Course
	if err := decodeJSONBody(r, &course); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if validationErr := validateCourse(course); validationErr != nil {
		writeError(w, http.StatusBadRequest, validationErr.Error())
		return
	}

	if err := h.repo.Create(r.Context(), &course); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create course")
		return
	}

	writeJSON(w, http.StatusCreated, course)
}

func (h *CoursesHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid course id")
		return
	}

	var course models.Course
	if err := decodeJSONBody(r, &course); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	course.ID = id

	if validationErr := validateCourse(course); validationErr != nil {
		writeError(w, http.StatusBadRequest, validationErr.Error())
		return
	}

	err = h.repo.Update(r.Context(), course)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "course not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update course")
		return
	}

	writeJSON(w, http.StatusOK, course)
}

func (h *CoursesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid course id")
		return
	}

	err = h.repo.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "course not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete course")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "course deleted"})
}

func validateCourse(course models.Course) error {
	if strings.TrimSpace(course.Title) == "" {
		return errors.New("title is required")
	}
	if course.TeacherID <= 0 {
		return errors.New("teacher_id must be positive")
	}

	return nil
}
