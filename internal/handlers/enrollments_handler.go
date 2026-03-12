package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"go-uni/internal/models"
	"go-uni/internal/repository"
	"go-uni/pkg/middleware"
)

type EnrollmentsHandler struct {
	repo *repository.EnrollmentsRepository
}

func NewEnrollmentsHandler(repo *repository.EnrollmentsRepository) *EnrollmentsHandler {
	return &EnrollmentsHandler{repo: repo}
}

// Enroll godoc
// @Summary Enroll student to course
// @Description Creates student-course enrollment.
// @Tags enrollments
// @Produce json
// @Param id path int true "Student ID"
// @Param course_id path int true "Course ID"
// @Security BearerAuth
// @Success 201 {object} models.Enrollment
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /students/{id}/courses/{course_id} [post]
func (h *EnrollmentsHandler) Enroll(w http.ResponseWriter, r *http.Request) {
	studentID, courseID, err := parseEnrollmentPathIDs(r)
	if err != nil {
		middleware.LogHandlerError(r, "invalid student/course path", err)
		_ = writeJSONError(w, http.StatusBadRequest, "invalid student/course path")
		return
	}

	enrollment := models.Enrollment{StudentID: studentID, CourseID: courseID}
	if err := h.repo.Enroll(r.Context(), enrollment); err != nil {
		middleware.LogHandlerError(r, "failed to enroll student", err)
		_ = writeJSONError(w, http.StatusInternalServerError, "failed to enroll student")
		return
	}

	jsonResponse(w, http.StatusCreated, enrollment)
}

// Unenroll godoc
// @Summary Unenroll student from course
// @Description Removes student-course enrollment.
// @Tags enrollments
// @Produce json
// @Param id path int true "Student ID"
// @Param course_id path int true "Course ID"
// @Security BearerAuth
// @Success 200 {object} messageResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /students/{id}/courses/{course_id} [delete]
func (h *EnrollmentsHandler) Unenroll(w http.ResponseWriter, r *http.Request) {
	studentID, courseID, err := parseEnrollmentPathIDs(r)
	if err != nil {
		middleware.LogHandlerError(r, "invalid student/course path", err)
		_ = writeJSONError(w, http.StatusBadRequest, "invalid student/course path")
		return
	}

	err = h.repo.Unenroll(r.Context(), studentID, courseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			middleware.LogHandlerError(r, "enrollment not found", err)
			_ = writeJSONError(w, http.StatusNotFound, "enrollment not found")
			return
		}
		middleware.LogHandlerError(r, "failed to unenroll student", err)
		_ = writeJSONError(w, http.StatusInternalServerError, "failed to unenroll student")
		return
	}

	jsonResponse(w, http.StatusOK, messageResponse{Message: "student unenrolled"})
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
