package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"go-uni/internal/models"
	"go-uni/internal/repository"
	"go-uni/pkg/middleware"
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
		middleware.LogHandlerError(r, "failed to list students", err)
		_ = writeJSONError(w, http.StatusInternalServerError, "failed to list students")
		return
	}

	jsonResponse(w, http.StatusOK, students)
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
		middleware.LogHandlerError(r, "invalid student id", err)
		_ = writeJSONError(w, http.StatusBadRequest, "invalid student id")
		return
	}

	student, getErr := h.repo.GetByID(r.Context(), id)
	if getErr != nil {
		middleware.LogHandlerError(r, "failed to get student", getErr)
		_ = writeJSONError(w, http.StatusInternalServerError, "failed to get student")
		return
	}
	if student == nil {
		middleware.LogHandlerError(r, "student not found", nil)
		_ = writeJSONError(w, http.StatusNotFound, "student not found")
		return
	}

	jsonResponse(w, http.StatusOK, student)
}

// Create godoc
// @Summary Create student
// @Description Creates a new student.
// @Tags students
// @Accept json
// @Produce json
// @Param request body models.CreateUpdateStudentRequest true "Student payload"
// @Security BearerAuth
// @Success 201 {object} models.Student
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /students [post]
func (h *StudentsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUpdateStudentRequest
	if err := readJSON(w, r, &req); err != nil {
		middleware.LogHandlerError(r, "invalid JSON body", err)
		_ = writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	student := models.Student{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	if validationErr := validateStudent(student); validationErr != nil {
		middleware.LogHandlerError(r, "student validation failed", validationErr)
		_ = writeJSONError(w, http.StatusBadRequest, validationErr.Error())
		return
	}

	if err := h.repo.Create(r.Context(), &student); err != nil {
		middleware.LogHandlerError(r, "failed to create student", err)
		_ = writeJSONError(w, http.StatusInternalServerError, "failed to create student")
		return
	}

	jsonResponse(w, http.StatusCreated, student)
}

// Update godoc
// @Summary Update student
// @Description Updates a student by ID.
// @Tags students
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Param request body models.CreateUpdateStudentRequest true "Student payload"
// @Security BearerAuth
// @Success 200 {object} models.Student
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /students/{id} [put]
func (h *StudentsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathID(r, "id")
	if err != nil {
		middleware.LogHandlerError(r, "invalid student id", err)
		_ = writeJSONError(w, http.StatusBadRequest, "invalid student id")
		return
	}

	var req models.CreateUpdateStudentRequest
	if err := readJSON(w, r, &req); err != nil {
		middleware.LogHandlerError(r, "invalid JSON body", err)
		_ = writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	student := models.Student{
		ID:        id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	if validationErr := validateStudent(student); validationErr != nil {
		middleware.LogHandlerError(r, "student validation failed", validationErr)
		_ = writeJSONError(w, http.StatusBadRequest, validationErr.Error())
		return
	}

	err = h.repo.Update(r.Context(), student)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			middleware.LogHandlerError(r, "student not found", err)
			_ = writeJSONError(w, http.StatusNotFound, "student not found")
			return
		}
		middleware.LogHandlerError(r, "failed to update student", err)
		_ = writeJSONError(w, http.StatusInternalServerError, "failed to update student")
		return
	}

	jsonResponse(w, http.StatusOK, student)
}

// Delete godoc
// @Summary Delete student
// @Description Deletes a student by ID.
// @Tags students
// @Produce json
// @Param id path int true "Student ID"
// @Security BearerAuth
// @Success 200 {object} messageResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /students/{id} [delete]
func (h *StudentsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathID(r, "id")
	if err != nil {
		middleware.LogHandlerError(r, "invalid student id", err)
		_ = writeJSONError(w, http.StatusBadRequest, "invalid student id")
		return
	}

	err = h.repo.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			middleware.LogHandlerError(r, "student not found", err)
			_ = writeJSONError(w, http.StatusNotFound, "student not found")
			return
		}
		middleware.LogHandlerError(r, "failed to delete student", err)
		_ = writeJSONError(w, http.StatusInternalServerError, "failed to delete student")
		return
	}

	jsonResponse(w, http.StatusOK, messageResponse{Message: "student deleted"})
}

func validateStudent(student models.Student) error {
	return validatePayload(student)
}
