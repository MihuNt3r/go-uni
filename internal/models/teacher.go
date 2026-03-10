package models

type Teacher struct {
	ID         int64  `json:"id"`
	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name" validate:"required"`
	Department string `json:"department" validate:"required"`
}
