package models

type User struct {
	Id        int     `json:"id,omitempty" db:"id"`
	Login     string  `json:"login" db:"login"`
	Email     string  `json:"email" db:"email"`
	Password  string  `json:"password,omitempty" db:"password"`
	FirstName *string `json:"firstName,omitempty" db:"first_name"`
	LastName  *string `json:"lastName,omitempty" db:"last_name"`
	Phone     *string `json:"phone,omitempty" db:"phone"`
}
