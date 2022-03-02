package service

import (
	"errors"
	"shop_backend/internal/models"
	"shop_backend/internal/repository"
	"shop_backend/pkg/auth"
	"shop_backend/pkg/hash"
	"time"
)

type UsersService struct {
	repo         repository.Users
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUsersService(repo repository.Users,
	hasher hash.PasswordHasher,
	tokenManager auth.TokenManager,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *UsersService {
	return &UsersService{
		repo:            repo,
		hasher:          hasher,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
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

func (s *UsersService) SignIn(email, password string) (models.Tokens, error) {
	passwordHash, err := s.hasher.Hash(password)
	if err != nil {
		return models.Tokens{}, err
	}

	user, err := s.repo.GetByCredentials(email, passwordHash)
	if err != nil {
		return models.Tokens{}, errors.New("user not found")
	}

	return s.createSession(user.Id)
}

func (s *UsersService) createSession(userId int) (models.Tokens, error) {
	var (
		res models.Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(userId, s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewJWT(userId, s.refreshTokenTTL)
	if err != nil {
		return res, err
	}

	return res, err
}

func (s *UsersService) GetUserById(id int) (models.User, error) {
	return s.repo.GetById(id)
}
