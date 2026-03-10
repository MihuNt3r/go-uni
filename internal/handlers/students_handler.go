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

func (h *StudentsHandler) List(w http.ResponseWriter, r *http.Request) {
	students, err := h.repo.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list students")
		return
	}

	writeJSON(w, http.StatusOK, students)
}

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

	writeJSON(w, http.StatusOK, map[string]string{"message": "student deleted"})
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
