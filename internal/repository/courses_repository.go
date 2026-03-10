package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go-uni/internal/models"
)

type CoursesRepository struct {
	db *sql.DB
}

func NewCoursesRepository(db *sql.DB) *CoursesRepository {
	return &CoursesRepository{db: db}
}

func (r *CoursesRepository) Create(ctx context.Context, course *models.Course) error {
	if course == nil {
		return errors.New("course is nil")
	}

	query := `
		INSERT INTO courses (title, description, teacher_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	if err := r.db.QueryRowContext(ctx, query, course.Title, course.Description, course.TeacherID).Scan(&course.ID); err != nil {
		return fmt.Errorf("create course: %w", err)
	}

	return nil
}

func (r *CoursesRepository) GetByID(ctx context.Context, id int64) (*models.Course, error) {
	query := `
		SELECT id, title, description, teacher_id
		FROM courses
		WHERE id = $1
	`

	var course models.Course
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&course.ID,
		&course.Title,
		&course.Description,
		&course.TeacherID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get course by id: %w", err)
	}

	return &course, nil
}

func (r *CoursesRepository) GetAll(ctx context.Context) ([]models.Course, error) {
	query := `
		SELECT id, title, description, teacher_id
		FROM courses
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list courses: %w", err)
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		if scanErr := rows.Scan(&course.ID, &course.Title, &course.Description, &course.TeacherID); scanErr != nil {
			return nil, fmt.Errorf("scan course row: %w", scanErr)
		}
		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate course rows: %w", err)
	}

	return courses, nil
}

func (r *CoursesRepository) Update(ctx context.Context, course models.Course) error {
	query := `
		UPDATE courses
		SET title = $1, description = $2, teacher_id = $3
		WHERE id = $4
	`

	res, err := r.db.ExecContext(ctx, query, course.Title, course.Description, course.TeacherID, course.ID)
	if err != nil {
		return fmt.Errorf("update course: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update course rows affected: %w", err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *CoursesRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM courses
		WHERE id = $1
	`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete course: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete course rows affected: %w", err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
