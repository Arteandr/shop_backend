package models

type User struct {
	id       string `db:"id"`
	email    string `db:"email"`
	password string `db:"password"`
}
