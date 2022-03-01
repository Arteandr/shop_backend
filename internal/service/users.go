package service

import (
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

func (s *UsersService) EmailExist(email string) bool {
	return s.repo.Exist(email)
}

func (s *UsersService) SignUp(email, password string) (int, error) {
	passwordHash, err := s.hasher.Hash(password)
	if err != nil {
		return 0, err
	}

	user := models.User{
		Email:    email,
		Password: passwordHash,
	}

	id, err := s.repo.Create(user)
	if err != nil {
		return 0, err
	}

	return id, nil
}
