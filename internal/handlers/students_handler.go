package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"go-uni/internal/models"
	"go-uni/internal/repository"
)

type StudentsHandler struct {
	repo *repository.StudentsRepository
}

func NewStudentsHandler(repo *repository.StudentsRepository) *StudentsHandler {
	return &StudentsHandler{repo: repo}
}

// List godoc
// @Summary List students
// @Description Returns all students.
// @Tags students
// @Produce json
// @Success 200 {array} models.Student
// @Failure 500 {object} errorResponse
// @Router /students [get]
func (h *StudentsHandler) List(w http.ResponseWriter, r *http.Request) {
	students, err := h.repo.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list students")
		return
	}

	writeJSON(w, http.StatusOK, students)
}

// Get godoc
// @Summary Get student by ID
// @Description Returns a single student by its ID.
// @Tags students
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} models.Student
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /students/{id} [get]
func (h *StudentsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid student id")
		return
	}

	student, getErr := h.repo.GetByID(r.Context(), id)
	if getErr != nil {
		writeError(w, http.StatusInternalServerError, "failed to get student")
		return
	}
	if student == nil {
		writeError(w, http.StatusNotFound, "student not found")
		return
	}

	writeJSON(w, http.StatusOK, student)
}

// Create godoc
// @Summary Create student
// @Description Creates a new student.
// @Tags students
// @Accept json
// @Produce json
// @Param request body models.Student true "Student payload"
// @Success 201 {object} models.Student
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /students [post]
func (h *StudentsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	if err := decodeJSONBody(r, &student); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if validationErr := validateStudent(student); validationErr != nil {
		writeError(w, http.StatusBadRequest, validationErr.Error())
		return
	}

	if err := h.repo.Create(r.Context(), &student); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create student")
		return
	}

	writeJSON(w, http.StatusCreated, student)
}

// Update godoc
// @Summary Update student
// @Description Updates a student by ID.
// @Tags students
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Param request body models.Student true "Student payload"
// @Success 200 {object} models.Student
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /students/{id} [put]
func (h *StudentsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid student id")
		return
	}

	var student models.Student
	if err := decodeJSONBody(r, &student); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	student.ID = id

	if validationErr := validateStudent(student); validationErr != nil {
		writeError(w, http.StatusBadRequest, validationErr.Error())
		return
	}

	err = h.repo.Update(r.Context(), student)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "student not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update student")
		return
	}

	writeJSON(w, http.StatusOK, student)
}

// Delete godoc
// @Summary Delete student
// @Description Deletes a student by ID.
// @Tags students
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} messageResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /students/{id} [delete]
func (h *StudentsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid student id")
		return
	}

	err = h.repo.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "student not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete student")
		return
	}

	writeJSON(w, http.StatusOK, messageResponse{Message: "student deleted"})
}

func validateStudent(student models.Student) error {
	if strings.TrimSpace(student.FirstName) == "" {
		return errors.New("first_name is required")
	}
	if strings.TrimSpace(student.LastName) == "" {
		return errors.New("last_name is required")
	}
	if strings.TrimSpace(student.Email) == "" {
		return errors.New("email is required")
	}

	return nil
}
