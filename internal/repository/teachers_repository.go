package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go-uni/internal/models"
)

type TeachersRepository struct {
	db *sql.DB
}

func NewTeachersRepository(db *sql.DB) *TeachersRepository {
	return &TeachersRepository{db: db}
}

func (r *TeachersRepository) Create(ctx context.Context, teacher *models.Teacher) error {
	if teacher == nil {
		return errors.New("teacher is nil")
	}

	query := `
		INSERT INTO teachers (first_name, last_name, department)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	if err := r.db.QueryRowContext(ctx, query, teacher.FirstName, teacher.LastName, teacher.Department).Scan(&teacher.ID); err != nil {
		return fmt.Errorf("create teacher: %w", err)
	}

	return nil
}

func (r *TeachersRepository) GetByID(ctx context.Context, id int64) (*models.Teacher, error) {
	query := `
		SELECT id, first_name, last_name, department
		FROM teachers
		WHERE id = $1
	`

	var teacher models.Teacher
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&teacher.ID,
		&teacher.FirstName,
		&teacher.LastName,
		&teacher.Department,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get teacher by id: %w", err)
	}

	return &teacher, nil
}

func (r *TeachersRepository) GetAll(ctx context.Context) ([]models.Teacher, error) {
	query := `
		SELECT id, first_name, last_name, department
		FROM teachers
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list teachers: %w", err)
	}
	defer rows.Close()

	var teachers []models.Teacher
	for rows.Next() {
		var teacher models.Teacher
		if scanErr := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Department); scanErr != nil {
			return nil, fmt.Errorf("scan teacher row: %w", scanErr)
		}
		teachers = append(teachers, teacher)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate teacher rows: %w", err)
	}

	return teachers, nil
}

func (r *TeachersRepository) Update(ctx context.Context, teacher models.Teacher) error {
	query := `
		UPDATE teachers
		SET first_name = $1, last_name = $2, department = $3
		WHERE id = $4
	`

	res, err := r.db.ExecContext(ctx, query, teacher.FirstName, teacher.LastName, teacher.Department, teacher.ID)
	if err != nil {
		return fmt.Errorf("update teacher: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update teacher rows affected: %w", err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *TeachersRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM teachers
		WHERE id = $1
	`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete teacher: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete teacher rows affected: %w", err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
