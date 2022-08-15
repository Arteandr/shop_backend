package models

type Address struct {
	Id      int    `json:"id,omitempty" db:"id"`
	Country string `json:"country" db:"country" binding:"required"`
	City    string `json:"city" db:"city" binding:"required"`
	Street  string `json:"street" db:"street" binding:"required"`
	Zip     string `json:"zip" db:"zip" binding:"required"`
}
