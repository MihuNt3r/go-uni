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

func (h *TeachersHandler) List(w http.ResponseWriter, r *http.Request) {
	teachers, err := h.repo.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list teachers")
		return
	}

	writeJSON(w, http.StatusOK, teachers)
}

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

	writeJSON(w, http.StatusOK, map[string]string{"message": "teacher deleted"})
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
