package repository

import (
	"context"
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

func (r *UsersRepo) Create(ctx context.Context, user models.User) (models.User, error) {
	var newUser models.User
	query := fmt.Sprintf("INSERT INTO %s (login, email, password) VALUES ($1,$2,$3) RETURNING *;", usersTable)
	if err := r.db.QueryRow(query, user.Login, user.Email, user.Password).Scan(&newUser.Id, &newUser.Email, &newUser.Login, &newUser.Password); err != nil {
		return models.User{}, err
	}

	return newUser, nil
}
