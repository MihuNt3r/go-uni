package models

type Course struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	TeacherID   int64  `json:"teacher_id"`
}
