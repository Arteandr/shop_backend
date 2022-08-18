package models

type Delivery struct {
	Id          int     `json:"id,omitempty" db:"id"`
	Name        string  `json:"name" db:"name"`
	CompanyName string  `json:"companyName" db:"company_name"`
	Price       float64 `json:"price" db:"price"`
}
