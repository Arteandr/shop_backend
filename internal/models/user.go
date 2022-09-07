package models

type User struct {
	Id              int      `json:"id,omitempty" db:"id"`
	Login           string   `json:"login" db:"login"`
	Email           string   `json:"email" db:"email"`
	Password        string   `json:"password,omitempty" db:"password"`
	FirstName       *string  `json:"firstName,omitempty" db:"first_name"`
	LastName        *string  `json:"lastName,omitempty" db:"last_name"`
	CompanyName     *string  `json:"companyName,omitempty" db:"company_name"`
	Phone           *Phone   `json:"phone,omitempty" db:"phone"`
	InvoiceAddress  *Address `json:"invoiceAddress,omitempty"`
	ShippingAddress *Address `json:"shippingAddress,omitempty"`
	Admin           bool     `json:"admin,omitempty" db:"admin"`
	Completed       bool     `json:"completed,omitempty" db:"completed"`
}
