package models

type Category struct {
	Id   int    `json:"id,omitempty" db:"id"`
	Name string `json:"name" binding:"required" db:"name"`
}
