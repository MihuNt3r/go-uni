package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go-uni/internal/models"
)

type StudentsRepository struct {
	db *sql.DB
}

func NewStudentsRepository(db *sql.DB) *StudentsRepository {
	return &StudentsRepository{db: db}
}

func (r *StudentsRepository) Create(ctx context.Context, student *models.Student) error {
	if student == nil {
		return errors.New("student is nil")
	}

	query := `
		INSERT INTO students (first_name, last_name, email)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	if err := r.db.QueryRowContext(ctx, query, student.FirstName, student.LastName, student.Email).Scan(&student.ID); err != nil {
		return fmt.Errorf("create student: %w", err)
	}

	return nil
}

func (r *StudentsRepository) GetByID(ctx context.Context, id int64) (*models.Student, error) {
	query := `
		SELECT id, first_name, last_name, email
		FROM students
		WHERE id = $1
	`

	var student models.Student
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&student.ID,
		&student.FirstName,
		&student.LastName,
		&student.Email,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get student by id: %w", err)
	}

	return &student, nil
}

func (r *StudentsRepository) GetAll(ctx context.Context) ([]models.Student, error) {
	query := `
		SELECT id, first_name, last_name, email
		FROM students
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list students: %w", err)
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var student models.Student
		if scanErr := rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email); scanErr != nil {
			return nil, fmt.Errorf("scan student row: %w", scanErr)
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate student rows: %w", err)
	}

	return students, nil
}

func (r *StudentsRepository) Update(ctx context.Context, student models.Student) error {
	query := `
		UPDATE students
		SET first_name = $1, last_name = $2, email = $3
		WHERE id = $4
	`

	res, err := r.db.ExecContext(ctx, query, student.FirstName, student.LastName, student.Email, student.ID)
	if err != nil {
		return fmt.Errorf("update student: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update student rows affected: %w", err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *StudentsRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM students
		WHERE id = $1
	`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete student: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete student rows affected: %w", err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
