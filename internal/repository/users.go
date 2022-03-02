package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"shop_backend/internal/models"
)

type UsersRepo struct {
	db *sqlx.DB
}

func NewUsersRepo(db *sqlx.DB) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (r *UsersRepo) Create(user models.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (email, password) VALUES ($1, $2) RETURNING id;", usersTable)
	row := r.db.QueryRow(query, user.Email, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UsersRepo) Exist(email string) bool {
	var user models.User
	query := fmt.Sprintf("SELECT * from %s where email=$1;", usersTable)
	err := r.db.Get(&user, query, email)
	if err != nil {
		return false
	}

	return true
}

func (r *UsersRepo) GetByCredentials(email, passwordHash string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT id, email FROM %s WHERE password=$1 and email=$2", usersTable)
	err := r.db.Get(&user, query, passwordHash, email)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
