package service

import (
	"context"
	"shop_backend/internal/models"
	"shop_backend/internal/repository"
	"shop_backend/pkg/auth"
	"shop_backend/pkg/hash"
)

type UsersService struct {
	repo         repository.Users
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
}

func NewUsersService(repo repository.Users, hasher hash.PasswordHasher, tokenManager auth.TokenManager) *UsersService {
	return &UsersService{
		repo:         repo,
		hasher:       hasher,
		tokenManager: tokenManager,
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

	newUser.Password = ""

	return newUser, err
}
