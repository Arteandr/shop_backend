package models

type Address struct {
	Id      int    `json:"id,omitempty" db:"id"`
	Country string `json:"country" db:"country"`
	City    string `json:"city" db:"city"`
	Street  string `json:"street" db:"street"`
	Zip     int    `json:"zip" db:"zip"`
}
