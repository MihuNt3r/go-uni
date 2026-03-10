package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"go-uni/internal/models"
	"go-uni/internal/repository"
)

type TeachersHandler struct {
	repo *repository.TeachersRepository
}

func NewTeachersHandler(repo *repository.TeachersRepository) *TeachersHandler {
	return &TeachersHandler{repo: repo}
}

// List godoc
// @Summary List teachers
// @Description Returns all teachers.
// @Tags teachers
// @Produce json
// @Success 200 {array} models.Teacher
// @Failure 500 {object} errorResponse
// @Router /teachers [get]
func (h *TeachersHandler) List(w http.ResponseWriter, r *http.Request) {
	teachers, err := h.repo.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list teachers")
		return
	}

	writeJSON(w, http.StatusOK, teachers)
}

// Get godoc
// @Summary Get teacher by ID
// @Description Returns a single teacher by its ID.
// @Tags teachers
// @Produce json
// @Param id path int true "Teacher ID"
// @Success 200 {object} models.Teacher
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /teachers/{id} [get]
func (h *TeachersHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid teacher id")
		return
	}

	teacher, getErr := h.repo.GetByID(r.Context(), id)
	if getErr != nil {
		writeError(w, http.StatusInternalServerError, "failed to get teacher")
		return
	}
	if teacher == nil {
		writeError(w, http.StatusNotFound, "teacher not found")
		return
	}

	writeJSON(w, http.StatusOK, teacher)
}

// Create godoc
// @Summary Create teacher
// @Description Creates a new teacher.
// @Tags teachers
// @Accept json
// @Produce json
// @Param request body models.Teacher true "Teacher payload"
// @Success 201 {object} models.Teacher
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /teachers [post]
func (h *TeachersHandler) Create(w http.ResponseWriter, r *http.Request) {
	var teacher models.Teacher
	if err := decodeJSONBody(r, &teacher); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if validationErr := validateTeacher(teacher); validationErr != nil {
		writeError(w, http.StatusBadRequest, validationErr.Error())
		return
	}

	if err := h.repo.Create(r.Context(), &teacher); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create teacher")
		return
	}

	writeJSON(w, http.StatusCreated, teacher)
}

// Update godoc
// @Summary Update teacher
// @Description Updates a teacher by ID.
// @Tags teachers
// @Accept json
// @Produce json
// @Param id path int true "Teacher ID"
// @Param request body models.Teacher true "Teacher payload"
// @Success 200 {object} models.Teacher
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /teachers/{id} [put]
func (h *TeachersHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid teacher id")
		return
	}

	var teacher models.Teacher
	if err := decodeJSONBody(r, &teacher); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	teacher.ID = id

	if validationErr := validateTeacher(teacher); validationErr != nil {
		writeError(w, http.StatusBadRequest, validationErr.Error())
		return
	}

	err = h.repo.Update(r.Context(), teacher)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "teacher not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update teacher")
		return
	}

	writeJSON(w, http.StatusOK, teacher)
}

// Delete godoc
// @Summary Delete teacher
// @Description Deletes a teacher by ID.
// @Tags teachers
// @Produce json
// @Param id path int true "Teacher ID"
// @Success 200 {object} messageResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /teachers/{id} [delete]
func (h *TeachersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid teacher id")
		return
	}

	err = h.repo.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "teacher not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete teacher")
		return
	}

	writeJSON(w, http.StatusOK, messageResponse{Message: "teacher deleted"})
}

func validateTeacher(teacher models.Teacher) error {
	if strings.TrimSpace(teacher.FirstName) == "" {
		return errors.New("first_name is required")
	}
	if strings.TrimSpace(teacher.LastName) == "" {
		return errors.New("last_name is required")
	}
	if strings.TrimSpace(teacher.Department) == "" {
		return errors.New("department is required")
	}

	return nil
}
