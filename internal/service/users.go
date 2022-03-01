package service

import (
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
