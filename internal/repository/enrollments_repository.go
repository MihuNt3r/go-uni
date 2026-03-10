package repository

import (
	"context"
	"database/sql"
	"fmt"

	"go-uni/internal/models"
)

type EnrollmentsRepository struct {
	db *sql.DB
}

func (r *EnrollmentsRepository) Enroll(ctx context.Context, enrollment models.Enrollment) error {
	query := `
		INSERT INTO enrollments (student_id, course_id)
		VALUES ($1, $2)
	`

	if _, err := r.db.ExecContext(ctx, query, enrollment.StudentID, enrollment.CourseID); err != nil {
		return fmt.Errorf("enroll student to course: %w", err)
	}

	return nil
}

func (r *EnrollmentsRepository) Unenroll(ctx context.Context, studentID, courseID int64) error {
	query := `
		DELETE FROM enrollments
		WHERE student_id = $1 AND course_id = $2
	`

	res, err := r.db.ExecContext(ctx, query, studentID, courseID)
	if err != nil {
		return fmt.Errorf("unenroll student from course: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("unenroll rows affected: %w", err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
