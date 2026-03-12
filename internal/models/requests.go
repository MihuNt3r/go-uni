package models

type CreateUpdateCourseRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	TeacherID   int64  `json:"teacher_id" validate:"gt=0"`
}

type CreateUpdateStudentRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
}

type CreateUpdateTeacherRequest struct {
	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name" validate:"required"`
	Department string `json:"department" validate:"required"`
}
