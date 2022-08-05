package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"shop_backend/internal/models"
	"time"
)

type UsersRepo struct {
	db *sqlx.DB
}

func NewUsersRepo(db *sqlx.DB) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

// $1 = login
// $2 = email
// $3 = password
func (r *UsersRepo) Create(ctx context.Context, user models.User) (models.User, error) {
	var newUser models.User
	query := fmt.Sprintf("INSERT INTO %s (login, email, password) VALUES ($1,$2,$3) RETURNING id, email, login;", usersTable)
	if err := r.db.QueryRow(query, user.Login, user.Email, user.Password).Scan(&newUser.Id, &newUser.Email, &newUser.Login); err != nil {
		return models.User{}, err
	}

	return newUser, nil
}

// $1 = login
// $2 = password
func (r *UsersRepo) GetByCredentials(ctx context.Context, findBy, login, password string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s=$1 AND password=$2;", usersTable, findBy)
	rows, err := r.db.Queryx(query, login, password)
	if err == sql.ErrNoRows {
		return models.User{}, models.ErrUserNotFound
	} else if err != nil {
		return models.User{}, err
	}

	for rows.Next() {
		if err := rows.StructScan(&user); err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}

// $1 = refreshToken
// $2 = time.Now()
func (r *UsersRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT U.* FROM %s AS S, %s AS U WHERE S.refresh_token=$1 AND S.expires_at > $2::timestamp AND U.id=S.user_id;", sessionsTable, usersTable)
	rows, err := r.db.Queryx(query, refreshToken, time.Now())
	if err == sql.ErrNoRows {
		return models.User{}, models.ErrUserNotFound
	} else if err != nil {
		return models.User{}, err
	}

	for rows.Next() {
		if err := rows.StructScan(&user); err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}

// $1 = userId
func (r *UsersRepo) GetById(ctx context.Context, userId int) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1;", usersTable)
	rows, err := r.db.Queryx(query, userId)
	if err == sql.ErrNoRows {
		return models.User{}, models.ErrUserNotFound
	} else if err != nil {
		return models.User{}, err
	}

	for rows.Next() {
		if err := rows.StructScan(&user); err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}

// $1 = userId
// $2 = refreshToken
// $3 = expiresAt
func (r *UsersRepo) SetSession(ctx context.Context, userId int, session models.Session) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id,refresh_token,expires_at) VALUES ($1,$2,$3);", sessionsTable)
	_, err := r.db.Exec(query, userId, session.RefreshToken, session.ExpiresAt)

	return err
}
