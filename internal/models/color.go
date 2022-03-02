package models

type Color struct {
	Id    int     `json:"id" db:"id"`
	Name  string  `json:"name" db:"name"`
	Hex   string  `json:"hex" db:"name"`
	Price float64 `json:"price" db:"price"`
}
