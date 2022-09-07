package models

type PaymentMethod struct {
	Id          int    `json:"id,omitempty" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Logo        string `json:"logo,omitempty"`
}
