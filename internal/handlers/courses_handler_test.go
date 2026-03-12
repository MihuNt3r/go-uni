package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-uni/internal/repository"
)

func newCoursesTestMux(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *http.ServeMux) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := repository.NewCoursesRepository(db)
	h := NewCoursesHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /courses", h.List)
	mux.HandleFunc("GET /courses/{id}", h.Get)
	mux.HandleFunc("POST /courses", h.Create)
	mux.HandleFunc("PUT /courses/{id}", h.Update)
	mux.HandleFunc("DELETE /courses/{id}", h.Delete)

	return db, mock, mux
}

func TestCoursesHandler_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, mux := newCoursesTestMux(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "title", "description", "teacher_id"}).
			AddRow(int64(1), "Math", "Algebra", int64(2)).
			AddRow(int64(2), "Physics", "Mechanics", int64(3))

		mock.ExpectQuery("FROM courses").WillReturnRows(rows)

		rr := executeRequest(mux, http.MethodGet, "/courses", "")

		require.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.JSONEq(t, `{"data":[
			{"id":1,"title":"Math","description":"Algebra","teacher_id":2},
			{"id":2,"title":"Physics","description":"Mechanics","teacher_id":3}
		]}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, mux := newCoursesTestMux(t)
		defer db.Close()

		mock.ExpectQuery("FROM courses").WillReturnError(errors.New("db down"))

		rr := executeRequest(mux, http.MethodGet, "/courses", "")

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error":"failed to list courses"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCoursesHandler_Get(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		_, _, mux := newCoursesTestMux(t)

		rr := executeRequest(mux, http.MethodGet, "/courses/abc", "")

		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"invalid course id"}`, rr.Body.String())
	})

	t.Run("not found", func(t *testing.T) {
		db, mock, mux := newCoursesTestMux(t)
		defer db.Close()

		mock.ExpectQuery("FROM courses").WithArgs(int64(99)).WillReturnError(sql.ErrNoRows)

		rr := executeRequest(mux, http.MethodGet, "/courses/99", "")

		require.Equal(t, http.StatusNotFound, rr.Code)
		assert.JSONEq(t, `{"error":"course not found"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, mux := newCoursesTestMux(t)
		defer db.Close()

		mock.ExpectQuery("FROM courses").WithArgs(int64(1)).WillReturnError(errors.New("boom"))

		rr := executeRequest(mux, http.MethodGet, "/courses/1", "")

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error":"failed to get course"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success", func(t *testing.T) {
		db, mock, mux := newCoursesTestMux(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "title", "description", "teacher_id"}).
			AddRow(int64(1), "Math", "Algebra", int64(2))

		mock.ExpectQuery("FROM courses").WithArgs(int64(1)).WillReturnRows(rows)

		rr := executeRequest(mux, http.MethodGet, "/courses/1", "")

		require.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, `{"data":{"id":1,"title":"Math","description":"Algebra","teacher_id":2}}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCoursesHandler_Create(t *testing.T) {
	t.Run("invalid JSON", func(t *testing.T) {
		_, _, mux := newCoursesTestMux(t)

		rr := executeRequest(mux, http.MethodPost, "/courses", `{"title":"Math","teacher_id":1,"unknown":true}`)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"invalid JSON body"}`, rr.Body.String())
	})

	t.Run("validation error", func(t *testing.T) {
		_, _, mux := newCoursesTestMux(t)

		rr := executeRequest(mux, http.MethodPost, "/courses", `{"title":"Math","description":"Algebra","teacher_id":0}`)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"teacher_id must be positive"}`, rr.Body.String())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, mux := newCoursesTestMux(t)
		defer db.Close()

		mock.ExpectQuery("INSERT INTO courses").
			WithArgs("Math", "Algebra", int64(2)).
			WillReturnError(errors.New("insert failed"))

		rr := executeRequest(mux, http.MethodPost, "/courses", `{"title":"Math","description":"Algebra","teacher_id":2}`)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error":"failed to create course"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success", func(t *testing.T) {
		db, mock, mux := newCoursesTestMux(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id"}).AddRow(int64(11))
		mock.ExpectQuery("INSERT INTO courses").
			WithArgs("Math", "Algebra", int64(2)).
			WillReturnRows(rows)

		rr := executeRequest(mux, http.MethodPost, "/courses", `{"title":"Math","description":"Algebra","teacher_id":2}`)

		require.Equal(t, http.StatusCreated, rr.Code)
		assert.JSONEq(t, `{"data":{"id":11,"title":"Math","description":"Algebra","teacher_id":2}}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCoursesHandler_Update(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		_, _, mux := newCoursesTestMux(t)

		rr := executeRequest(mux, http.MethodPut, "/courses/abc", `{"title":"Math","description":"Algebra","teacher_id":2}`)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"invalid course id"}`, rr.Body.String())
	})

	t.Run("invalid JSON", func(t *testing.T) {
		_, _, mux := newCoursesTestMux(t)

		rr := executeRequest(mux, http.MethodPut, "/courses/1", `{"title":"Math","teacher_id":2,"extra":1}`)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"invalid JSON body"}`, rr.Body.String())
	})

	t.Run("validation error", func(t *testing.T) {
		_, _, mux := newCoursesTestMux(t)

		rr := executeRequest(mux, http.MethodPut, "/courses/1", `{"title":"","description":"Algebra","teacher_id":2}`)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"title is required"}`, rr.Body.String())
	})

	t.Run("not found", func(t *testing.T) {
		db, mock, mux := newCoursesTestMux(t)
		defer db.Close()

		mock.ExpectExec("UPDATE courses").
			WithArgs("Math", "Algebra", int64(2), int64(99)).
			WillReturnResult(sqlmock.NewResult(0, 0))

		rr := executeRequest(mux, http.MethodPut, "/courses/99", `{"title":"Math","description":"Algebra","teacher_id":2}`)

		require.Equal(t, http.StatusNotFound, rr.Code)
		assert.JSONEq(t, `{"error":"course not found"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, mux := newCoursesTestMux(t)
		defer db.Close()

		mock.ExpectExec("UPDATE courses").
			WithArgs("Math", "Algebra", int64(2), int64(1)).
			WillReturnError(errors.New("update failed"))

		rr := executeRequest(mux, http.MethodPut, "/courses/1", `{"title":"Math","description":"Algebra","teacher_id":2}`)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error":"failed to update course"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success", func(t *testing.T) {
		db, mock, mux := newCoursesTestMux(t)
		defer db.Close()

		mock.ExpectExec("UPDATE courses").
			WithArgs("Math", "Algebra", int64(2), int64(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		rr := executeRequest(mux, http.MethodPut, "/courses/1", `{"title":"Math","description":"Algebra","teacher_id":2}`)

		require.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, `{"data":{"id":1,"title":"Math","description":"Algebra","teacher_id":2}}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCoursesHandler_Delete(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		_, _, mux := newCoursesTestMux(t)

		rr := executeRequest(mux, http.MethodDelete, "/courses/abc", "")

		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"invalid course id"}`, rr.Body.String())
	})

	t.Run("not found", func(t *testing.T) {
		db, mock, mux := newCoursesTestMux(t)
		defer db.Close()

		mock.ExpectExec("DELETE FROM courses").WithArgs(int64(42)).WillReturnResult(sqlmock.NewResult(0, 0))

		rr := executeRequest(mux, http.MethodDelete, "/courses/42", "")

		require.Equal(t, http.StatusNotFound, rr.Code)
		assert.JSONEq(t, `{"error":"course not found"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, mux := newCoursesTestMux(t)
		defer db.Close()

		mock.ExpectExec("DELETE FROM courses").WithArgs(int64(1)).WillReturnError(errors.New("delete failed"))

		rr := executeRequest(mux, http.MethodDelete, "/courses/1", "")

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error":"failed to delete course"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success", func(t *testing.T) {
		db, mock, mux := newCoursesTestMux(t)
		defer db.Close()

		mock.ExpectExec("DELETE FROM courses").WithArgs(int64(1)).WillReturnResult(sqlmock.NewResult(0, 1))

		rr := executeRequest(mux, http.MethodDelete, "/courses/1", "")

		require.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, `{"data": {"message":"course deleted"}}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
