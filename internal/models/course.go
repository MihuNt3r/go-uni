package models

type Course struct {
	ID          int64  `json:"id"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	TeacherID   int64  `json:"teacher_id" validate:"gt=0"`
}
