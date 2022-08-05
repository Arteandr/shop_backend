package repository

import (
	"context"
	"database/sql"
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

func (r *UsersRepo) GetByCredentials(ctx context.Context, findBy, login, password string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s=$1 AND password=$2;", usersTable, findBy)
	if err := r.db.QueryRow(query, login, password).Scan(&user.Id, &user.Email, &user.Login, &user.Password); err == sql.ErrNoRows {
		return models.User{}, models.ErrUserNotFound
	} else if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *UsersRepo) SetSession(ctx context.Context, userId int, session models.Session) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id,refresh_token,expires_at) VALUES ($1,$2,$3);", sessionsTable)
	_, err := r.db.Exec(query, userId, session.RefreshToken, session.ExpiresAt)

	return err
}
