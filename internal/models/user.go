package models

type User struct {
	Id       int    `json:"id,omitempty" db:"id"`
	Login    string `json:"login" db:"login"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password,omitempty" db:"password"`
}
