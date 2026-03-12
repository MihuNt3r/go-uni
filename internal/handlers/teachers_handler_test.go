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

func newTeachersTestMux(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *http.ServeMux) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := repository.NewTeachersRepository(db)
	h := NewTeachersHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /teachers", h.List)
	mux.HandleFunc("GET /teachers/{id}", h.Get)
	mux.HandleFunc("POST /teachers", h.Create)
	mux.HandleFunc("PUT /teachers/{id}", h.Update)
	mux.HandleFunc("DELETE /teachers/{id}", h.Delete)

	return db, mock, mux
}

func TestTeachersHandler_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, mux := newTeachersTestMux(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "department"}).
			AddRow(int64(1), "John", "Doe", "Math").
			AddRow(int64(2), "Jane", "Smith", "Physics")

		mock.ExpectQuery("FROM teachers").WillReturnRows(rows)

		rr := executeRequest(mux, http.MethodGet, "/teachers", "")

		require.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.JSONEq(t, `{"data":[
			{"id":1,"first_name":"John","last_name":"Doe","department":"Math"},
			{"id":2,"first_name":"Jane","last_name":"Smith","department":"Physics"}
		]}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, mux := newTeachersTestMux(t)
		defer db.Close()

		mock.ExpectQuery("FROM teachers").WillReturnError(errors.New("db down"))

		rr := executeRequest(mux, http.MethodGet, "/teachers", "")

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error":"failed to list teachers"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestTeachersHandler_Get(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		_, _, mux := newTeachersTestMux(t)

		rr := executeRequest(mux, http.MethodGet, "/teachers/abc", "")

		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"invalid teacher id"}`, rr.Body.String())
	})

	t.Run("not found", func(t *testing.T) {
		db, mock, mux := newTeachersTestMux(t)
		defer db.Close()

		mock.ExpectQuery("FROM teachers").WithArgs(int64(99)).WillReturnError(sql.ErrNoRows)

		rr := executeRequest(mux, http.MethodGet, "/teachers/99", "")

		require.Equal(t, http.StatusNotFound, rr.Code)
		assert.JSONEq(t, `{"error":"teacher not found"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, mux := newTeachersTestMux(t)
		defer db.Close()

		mock.ExpectQuery("FROM teachers").WithArgs(int64(1)).WillReturnError(errors.New("boom"))

		rr := executeRequest(mux, http.MethodGet, "/teachers/1", "")

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error":"failed to get teacher"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success", func(t *testing.T) {
		db, mock, mux := newTeachersTestMux(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "department"}).
			AddRow(int64(1), "John", "Doe", "Math")

		mock.ExpectQuery("FROM teachers").WithArgs(int64(1)).WillReturnRows(rows)

		rr := executeRequest(mux, http.MethodGet, "/teachers/1", "")

		require.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, `{"data":{"id":1,"first_name":"John","last_name":"Doe","department":"Math"}}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestTeachersHandler_Create(t *testing.T) {
	t.Run("invalid JSON", func(t *testing.T) {
		_, _, mux := newTeachersTestMux(t)

		rr := executeRequest(mux, http.MethodPost, "/teachers", `{"first_name":"John","last_name":"Doe","department":"Math","unknown":true}`)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"invalid JSON body"}`, rr.Body.String())
	})

	t.Run("validation error", func(t *testing.T) {
		_, _, mux := newTeachersTestMux(t)

		rr := executeRequest(mux, http.MethodPost, "/teachers", `{"first_name":"John","last_name":"Doe","department":""}`)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"department is required"}`, rr.Body.String())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, mux := newTeachersTestMux(t)
		defer db.Close()

		mock.ExpectQuery("INSERT INTO teachers").
			WithArgs("John", "Doe", "Math").
			WillReturnError(errors.New("insert failed"))

		rr := executeRequest(mux, http.MethodPost, "/teachers", `{"first_name":"John","last_name":"Doe","department":"Math"}`)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error":"failed to create teacher"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success", func(t *testing.T) {
		db, mock, mux := newTeachersTestMux(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id"}).AddRow(int64(11))
		mock.ExpectQuery("INSERT INTO teachers").
			WithArgs("John", "Doe", "Math").
			WillReturnRows(rows)

		rr := executeRequest(mux, http.MethodPost, "/teachers", `{"first_name":"John","last_name":"Doe","department":"Math"}`)

		require.Equal(t, http.StatusCreated, rr.Code)
		assert.JSONEq(t, `{"data":{"id":11,"first_name":"John","last_name":"Doe","department":"Math"}}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestTeachersHandler_Update(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		_, _, mux := newTeachersTestMux(t)

		rr := executeRequest(mux, http.MethodPut, "/teachers/abc", `{"first_name":"John","last_name":"Doe","department":"Math"}`)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"invalid teacher id"}`, rr.Body.String())
	})

	t.Run("invalid JSON", func(t *testing.T) {
		_, _, mux := newTeachersTestMux(t)

		rr := executeRequest(mux, http.MethodPut, "/teachers/1", `{"first_name":"John","last_name":"Doe","department":"Math","extra":1}`)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"invalid JSON body"}`, rr.Body.String())
	})

	t.Run("validation error", func(t *testing.T) {
		_, _, mux := newTeachersTestMux(t)

		rr := executeRequest(mux, http.MethodPut, "/teachers/1", `{"first_name":"John","last_name":"Doe","department":""}`)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"department is required"}`, rr.Body.String())
	})

	t.Run("not found", func(t *testing.T) {
		db, mock, mux := newTeachersTestMux(t)
		defer db.Close()

		mock.ExpectExec("UPDATE teachers").
			WithArgs("John", "Doe", "Math", int64(99)).
			WillReturnResult(sqlmock.NewResult(0, 0))

		rr := executeRequest(mux, http.MethodPut, "/teachers/99", `{"first_name":"John","last_name":"Doe","department":"Math"}`)

		require.Equal(t, http.StatusNotFound, rr.Code)
		assert.JSONEq(t, `{"error":"teacher not found"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, mux := newTeachersTestMux(t)
		defer db.Close()

		mock.ExpectExec("UPDATE teachers").
			WithArgs("John", "Doe", "Math", int64(1)).
			WillReturnError(errors.New("update failed"))

		rr := executeRequest(mux, http.MethodPut, "/teachers/1", `{"first_name":"John","last_name":"Doe","department":"Math"}`)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error":"failed to update teacher"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success", func(t *testing.T) {
		db, mock, mux := newTeachersTestMux(t)
		defer db.Close()

		mock.ExpectExec("UPDATE teachers").
			WithArgs("John", "Doe", "Math", int64(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		rr := executeRequest(mux, http.MethodPut, "/teachers/1", `{"first_name":"John","last_name":"Doe","department":"Math"}`)

		require.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, `{"data":{"id":1,"first_name":"John","last_name":"Doe","department":"Math"}}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestTeachersHandler_Delete(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		_, _, mux := newTeachersTestMux(t)

		rr := executeRequest(mux, http.MethodDelete, "/teachers/abc", "")

		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"invalid teacher id"}`, rr.Body.String())
	})

	t.Run("not found", func(t *testing.T) {
		db, mock, mux := newTeachersTestMux(t)
		defer db.Close()

		mock.ExpectExec("DELETE FROM teachers").WithArgs(int64(42)).WillReturnResult(sqlmock.NewResult(0, 0))

		rr := executeRequest(mux, http.MethodDelete, "/teachers/42", "")

		require.Equal(t, http.StatusNotFound, rr.Code)
		assert.JSONEq(t, `{"error":"teacher not found"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, mux := newTeachersTestMux(t)
		defer db.Close()

		mock.ExpectExec("DELETE FROM teachers").WithArgs(int64(1)).WillReturnError(errors.New("delete failed"))

		rr := executeRequest(mux, http.MethodDelete, "/teachers/1", "")

		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error":"failed to delete teacher"}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success", func(t *testing.T) {
		db, mock, mux := newTeachersTestMux(t)
		defer db.Close()

		mock.ExpectExec("DELETE FROM teachers").WithArgs(int64(1)).WillReturnResult(sqlmock.NewResult(0, 1))

		rr := executeRequest(mux, http.MethodDelete, "/teachers/1", "")

		require.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, `{"data":{"message":"teacher deleted"}}`, rr.Body.String())
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
