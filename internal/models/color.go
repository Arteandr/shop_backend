package models

type Color struct {
	Id    int     `json:"id,omitempty" db:"id"`
	Name  string  `json:"name" binding:"required" db:"name"`
	Hex   string  `json:"hex" binding:"required" db:"name"`
	Price float64 `json:"price" binding:"required" db:"price"`
}
