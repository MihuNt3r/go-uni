package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"go-uni/internal/models"
	"go-uni/internal/repository"
)

type CoursesHandler struct {
	repo *repository.CoursesRepository
}

func NewCoursesHandler(repo *repository.CoursesRepository) *CoursesHandler {
	return &CoursesHandler{repo: repo}
}

// List godoc
// @Summary List courses
// @Description Returns all courses.
// @Tags courses
// @Produce json
// @Success 200 {array} models.Course
// @Failure 500 {object} errorResponse
// @Router /courses [get]
func (h *CoursesHandler) List(w http.ResponseWriter, r *http.Request) {
	courses, err := h.repo.GetAll(r.Context())
	if err != nil {
		_ = writeJSONError(w, http.StatusInternalServerError, "failed to list courses")
		return
	}

	jsonResponse(w, http.StatusOK, courses)
}

// Get godoc
// @Summary Get course by ID
// @Description Returns a single course by its ID.
// @Tags courses
// @Produce json
// @Param id path int true "Course ID"
// @Success 200 {object} models.Course
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /courses/{id} [get]
func (h *CoursesHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathID(r, "id")
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, "invalid course id")
		return
	}

	course, getErr := h.repo.GetByID(r.Context(), id)
	if getErr != nil {
		_ = writeJSONError(w, http.StatusInternalServerError, "failed to get course")
		return
	}
	if course == nil {
		_ = writeJSONError(w, http.StatusNotFound, "course not found")
		return
	}

	writeJSON(w, http.StatusOK, course)
}

// Create godoc
// @Summary Create course
// @Description Creates a new course.
// @Tags courses
// @Accept json
// @Produce json
// @Param request body models.CreateUpdateCourseRequest true "Course payload"
// @Success 201 {object} models.Course
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /courses [post]
func (h *CoursesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUpdateCourseRequest
	if err := readJSON(w, r, &req); err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	course := models.Course{
		Title:       req.Title,
		Description: req.Description,
		TeacherID:   req.TeacherID,
	}

	if validationErr := validateCourse(course); validationErr != nil {
		_ = writeJSONError(w, http.StatusBadRequest, validationErr.Error())
		return
	}

	if err := h.repo.Create(r.Context(), &course); err != nil {
		_ = writeJSONError(w, http.StatusInternalServerError, "failed to create course")
		return
	}

	jsonResponse(w, http.StatusCreated, course)
}

// Update godoc
// @Summary Update course
// @Description Updates a course by ID.
// @Tags courses
// @Accept json
// @Produce json
// @Param id path int true "Course ID"
// @Param request body models.CreateUpdateCourseRequest true "Course payload"
// @Success 200 {object} models.Course
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /courses/{id} [put]
func (h *CoursesHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathID(r, "id")
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, "invalid course id")
		return
	}

	var req models.CreateUpdateCourseRequest
	if err := readJSON(w, r, &req); err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	course := models.Course{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		TeacherID:   req.TeacherID,
	}

	if validationErr := validateCourse(course); validationErr != nil {
		_ = writeJSONError(w, http.StatusBadRequest, validationErr.Error())
		return
	}

	err = h.repo.Update(r.Context(), course)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = writeJSONError(w, http.StatusNotFound, "course not found")
			return
		}
		_ = writeJSONError(w, http.StatusInternalServerError, "failed to update course")
		return
	}

	jsonResponse(w, http.StatusOK, course)
}

// Delete godoc
// @Summary Delete course
// @Description Deletes a course by ID.
// @Tags courses
// @Produce json
// @Param id path int true "Course ID"
// @Success 200 {object} messageResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /courses/{id} [delete]
func (h *CoursesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathID(r, "id")
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, "invalid course id")
		return
	}

	err = h.repo.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = writeJSONError(w, http.StatusNotFound, "course not found")
			return
		}
		_ = writeJSONError(w, http.StatusInternalServerError, "failed to delete course")
		return
	}

	jsonResponse(w, http.StatusOK, messageResponse{Message: "course deleted"})
}

func validateCourse(course models.Course) error {
	return validatePayload(course)
}
