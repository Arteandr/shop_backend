package service

import (
	"context"
	"database/sql"
	"shop_backend/internal/models"
	"shop_backend/internal/repository"
	"shop_backend/pkg/auth"
	"shop_backend/pkg/hash"
	"strconv"
	"time"
)

type UsersService struct {
	repo         repository.Users
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUsersService(repo repository.Users, hasher hash.PasswordHasher, tokenManager auth.TokenManager, accessTokenTTL, refreshTokenTTL time.Duration) *UsersService {
	return &UsersService{
		repo:            repo,
		hasher:          hasher,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *UsersService) SignUp(ctx context.Context, email, login, password string) (models.User, error) {
	passwordHash, err := s.hasher.Hash(password)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		Email:    email,
		Login:    login,
		Password: passwordHash,
	}

	newUser, err := s.repo.Create(ctx, user)
	if err != nil {
		return models.User{}, err
	}

	// Hide password
	newUser.Password = ""

	return newUser, err
}

func (s *UsersService) SignIn(ctx context.Context, findBy, login, password string) (models.Tokens, error) {
	passwordHash, err := s.hasher.Hash(password)
	if err != nil {
		return models.Tokens{}, err
	}

	user, err := s.repo.GetByCredentials(ctx, findBy, login, passwordHash)
	if err != nil {
		return models.Tokens{}, err
	}

	return s.createSession(ctx, user.Id)
}

func (s *UsersService) Logout(ctx context.Context, userId int) error {
	if err := s.repo.DeleteSession(ctx, userId); err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

func (s *UsersService) createSession(ctx context.Context, userId int) (models.Tokens, error) {
	var (
		res models.Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(strconv.Itoa(userId), s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := models.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	err = s.repo.SetSession(ctx, userId, session)

	return res, err
}

func (s *UsersService) RefreshTokens(ctx context.Context, refreshToken string) (models.Tokens, error) {
	user, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return models.Tokens{}, err
	}

	return s.createSession(ctx, user.Id)
}

func (s *UsersService) GetById(ctx context.Context, userId int) (models.User, error) {
	user, err := s.repo.GetById(ctx, userId)
	if err != nil {
		return models.User{}, err
	}

	// Hide password
	user.Password = ""

	return user, nil
}
