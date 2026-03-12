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

func newStudentsTestMux(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *http.ServeMux) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := repository.NewStudentsRepository(db)
	h := NewStudentsHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /students", h.List)
	mux.HandleFunc("GET /students/{id}", h.Get)
	mux.HandleFunc("POST /students", h.Create)
	mux.HandleFunc("PUT /students/{id}", h.Update)
	mux.HandleFunc("DELETE /students/{id}", h.Delete)

	return db, mock, mux
}

func TestStudentsHandler_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, mux := newStudentsTestMux(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email"}).
			AddRow(int64(1), "John", "Doe", "john@example.com").
			AddRow(int64(2), "Jane", "Smith", "jane@example.com")

		mock.ExpectQuery("FROM students").WillReturnRows(rows)

		rr := executeRequest(mux, http.MethodGet, "/students", "")

		require.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.JSONEq(t, `[
			{"id":1,"first_name":"John","last_name":"Doe","email":"john@example.com"},
			{"id":2,"first_name":"Jane","last_name":"Smith","email":"jane@example.com"}
		]`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, mux := newStudentsTestMux(t)
		defer db.Close()

		mock.ExpectQuery("FROM students").WillReturnError(errors.New("db down"))

		rr := executeRequest(mux, http.MethodGet, "/students", "")

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error":"failed to list students"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestStudentsHandler_Get(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		_, _, mux := newStudentsTestMux(t)

		rr := executeRequest(mux, http.MethodGet, "/students/abc", "")

		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"invalid student id"}`, rr.Body.String())
	})

	t.Run("not found", func(t *testing.T) {
		db, mock, mux := newStudentsTestMux(t)
		defer db.Close()

		mock.ExpectQuery("FROM students").WithArgs(int64(99)).WillReturnError(sql.ErrNoRows)

		rr := executeRequest(mux, http.MethodGet, "/students/99", "")

		require.Equal(t, http.StatusNotFound, rr.Code)
		assert.JSONEq(t, `{"error":"student not found"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, mux := newStudentsTestMux(t)
		defer db.Close()

		mock.ExpectQuery("FROM students").WithArgs(int64(1)).WillReturnError(errors.New("boom"))

		rr := executeRequest(mux, http.MethodGet, "/students/1", "")

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error":"failed to get student"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success", func(t *testing.T) {
		db, mock, mux := newStudentsTestMux(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email"}).
			AddRow(int64(1), "John", "Doe", "john@example.com")

		mock.ExpectQuery("FROM students").WithArgs(int64(1)).WillReturnRows(rows)

		rr := executeRequest(mux, http.MethodGet, "/students/1", "")

		require.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, `{"id":1,"first_name":"John","last_name":"Doe","email":"john@example.com"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestStudentsHandler_Create(t *testing.T) {
	t.Run("invalid JSON", func(t *testing.T) {
		_, _, mux := newStudentsTestMux(t)

		rr := executeRequest(mux, http.MethodPost, "/students", `{"first_name":"John","last_name":"Doe","email":"john@example.com","unknown":true}`)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"invalid JSON body"}`, rr.Body.String())
	})

	t.Run("validation error", func(t *testing.T) {
		_, _, mux := newStudentsTestMux(t)

		rr := executeRequest(mux, http.MethodPost, "/students", `{"first_name":"John","last_name":"Doe","email":"not-an-email"}`)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"email must be valid"}`, rr.Body.String())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, mux := newStudentsTestMux(t)
		defer db.Close()

		mock.ExpectQuery("INSERT INTO students").
			WithArgs("John", "Doe", "john@example.com").
			WillReturnError(errors.New("insert failed"))

		rr := executeRequest(mux, http.MethodPost, "/students", `{"first_name":"John","last_name":"Doe","email":"john@example.com"}`)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error":"failed to create student"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success", func(t *testing.T) {
		db, mock, mux := newStudentsTestMux(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id"}).AddRow(int64(11))
		mock.ExpectQuery("INSERT INTO students").
			WithArgs("John", "Doe", "john@example.com").
			WillReturnRows(rows)

		rr := executeRequest(mux, http.MethodPost, "/students", `{"first_name":"John","last_name":"Doe","email":"john@example.com"}`)

		require.Equal(t, http.StatusCreated, rr.Code)
		assert.JSONEq(t, `{"id":11,"first_name":"John","last_name":"Doe","email":"john@example.com"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestStudentsHandler_Update(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		_, _, mux := newStudentsTestMux(t)

		rr := executeRequest(mux, http.MethodPut, "/students/abc", `{"first_name":"John","last_name":"Doe","email":"john@example.com"}`)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"invalid student id"}`, rr.Body.String())
	})

	t.Run("invalid JSON", func(t *testing.T) {
		_, _, mux := newStudentsTestMux(t)

		rr := executeRequest(mux, http.MethodPut, "/students/1", `{"first_name":"John","last_name":"Doe","email":"john@example.com","extra":1}`)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"invalid JSON body"}`, rr.Body.String())
	})

	t.Run("validation error", func(t *testing.T) {
		_, _, mux := newStudentsTestMux(t)

		rr := executeRequest(mux, http.MethodPut, "/students/1", `{"first_name":"John","last_name":"Doe","email":"invalid"}`)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"email must be valid"}`, rr.Body.String())
	})

	t.Run("not found", func(t *testing.T) {
		db, mock, mux := newStudentsTestMux(t)
		defer db.Close()

		mock.ExpectExec("UPDATE students").
			WithArgs("John", "Doe", "john@example.com", int64(99)).
			WillReturnResult(sqlmock.NewResult(0, 0))

		rr := executeRequest(mux, http.MethodPut, "/students/99", `{"first_name":"John","last_name":"Doe","email":"john@example.com"}`)

		require.Equal(t, http.StatusNotFound, rr.Code)
		assert.JSONEq(t, `{"error":"student not found"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, mux := newStudentsTestMux(t)
		defer db.Close()

		mock.ExpectExec("UPDATE students").
			WithArgs("John", "Doe", "john@example.com", int64(1)).
			WillReturnError(errors.New("update failed"))

		rr := executeRequest(mux, http.MethodPut, "/students/1", `{"first_name":"John","last_name":"Doe","email":"john@example.com"}`)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error":"failed to update student"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success", func(t *testing.T) {
		db, mock, mux := newStudentsTestMux(t)
		defer db.Close()

		mock.ExpectExec("UPDATE students").
			WithArgs("John", "Doe", "john@example.com", int64(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		rr := executeRequest(mux, http.MethodPut, "/students/1", `{"first_name":"John","last_name":"Doe","email":"john@example.com"}`)

		require.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, `{"id":1,"first_name":"John","last_name":"Doe","email":"john@example.com"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestStudentsHandler_Delete(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		_, _, mux := newStudentsTestMux(t)

		rr := executeRequest(mux, http.MethodDelete, "/students/abc", "")

		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"invalid student id"}`, rr.Body.String())
	})

	t.Run("not found", func(t *testing.T) {
		db, mock, mux := newStudentsTestMux(t)
		defer db.Close()

		mock.ExpectExec("DELETE FROM students").WithArgs(int64(42)).WillReturnResult(sqlmock.NewResult(0, 0))

		rr := executeRequest(mux, http.MethodDelete, "/students/42", "")

		require.Equal(t, http.StatusNotFound, rr.Code)
		assert.JSONEq(t, `{"error":"student not found"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, mux := newStudentsTestMux(t)
		defer db.Close()

		mock.ExpectExec("DELETE FROM students").WithArgs(int64(1)).WillReturnError(errors.New("delete failed"))

		rr := executeRequest(mux, http.MethodDelete, "/students/1", "")

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error":"failed to delete student"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success", func(t *testing.T) {
		db, mock, mux := newStudentsTestMux(t)
		defer db.Close()

		mock.ExpectExec("DELETE FROM students").WithArgs(int64(1)).WillReturnResult(sqlmock.NewResult(0, 1))

		rr := executeRequest(mux, http.MethodDelete, "/students/1", "")

		require.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, `{"message":"student deleted"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
