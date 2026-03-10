package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"go-uni/internal/models"
	"go-uni/internal/repository"
)

type EnrollmentsHandler struct {
	repo *repository.EnrollmentsRepository
}

func NewEnrollmentsHandler(repo *repository.EnrollmentsRepository) *EnrollmentsHandler {
	return &EnrollmentsHandler{repo: repo}
}

func (h *EnrollmentsHandler) Enroll(w http.ResponseWriter, r *http.Request) {
	studentID, courseID, err := parseEnrollmentPathIDs(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid student/course path")
		return
	}

	enrollment := models.Enrollment{StudentID: studentID, CourseID: courseID}
	if err := h.repo.Enroll(r.Context(), enrollment); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to enroll student")
		return
	}

	writeJSON(w, http.StatusCreated, enrollment)
}

func (h *EnrollmentsHandler) Unenroll(w http.ResponseWriter, r *http.Request) {
	studentID, courseID, err := parseEnrollmentPathIDs(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid student/course path")
		return
	}

	err = h.repo.Unenroll(r.Context(), studentID, courseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "enrollment not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to unenroll student")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "student unenrolled"})
}

func parseEnrollmentPathIDs(r *http.Request) (int64, int64, error) {
	studentID, studentErr := parsePathID(r, "id")
	if studentErr != nil {
		return 0, 0, studentErr
	}

	courseID, courseErr := parsePathID(r, "course_id")
	if courseErr != nil {
		return 0, 0, courseErr
	}

	return studentID, courseID, nil
}
