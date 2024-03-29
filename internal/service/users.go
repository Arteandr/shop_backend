package service

import (
	"context"
	"database/sql"
	"errors"
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

	if err := s.repo.CreatePhone(ctx, newUser.Id); err != nil {
		return models.User{}, err
	}

	if err := s.repo.CreateDefaultAddress(ctx, "invoice", newUser.Id); err != nil {
		return models.User{}, err

	}
	if err := s.repo.CreateDefaultAddress(ctx, "shipping", newUser.Id); err != nil {
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

func (s *UsersService) DeleteMe(ctx context.Context, userId int) error {
	return s.repo.Delete(ctx, userId)
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

func (s *UsersService) GetMe(ctx context.Context, userId int) (models.User, error) {
	user, err := s.repo.GetById(ctx, userId)
	if err != nil {
		return models.User{}, err
	}

	invoiceAddress, err := s.repo.GetAddress(ctx, "invoice", user.Id)
	if err != nil && !errors.Is(err, models.ErrAddressNotFound) {
		return models.User{}, err
	}
	if invoiceAddress != (models.Address{}) {
		user.InvoiceAddress = &invoiceAddress
	}

	shippingAddress, err := s.repo.GetAddress(ctx, "shipping", user.Id)
	if err != nil && !errors.Is(err, models.ErrAddressNotFound) {
		return models.User{}, err
	}
	if shippingAddress != (models.Address{}) {
		user.ShippingAddress = &shippingAddress
	}

	phone, err := s.repo.GetPhone(ctx, user.Id)
	if err != nil {
		return models.User{}, err
	}
	if phone.Code != nil && phone.Number != nil {
		phoneCode := *phone.Code
		phoneNumber := *phone.Number
		user.Phone = phoneCode + phoneNumber
	}

	// Hide password
	user.Password = ""

	return user, nil
}

func (s *UsersService) UpdateEmail(ctx context.Context, userId int, email string) error {
	return s.repo.UpdateField(ctx, "email", email, userId)
}

func (s *UsersService) UpdatePassword(ctx context.Context, userId int, oldPassword, newPassword string) error {
	user, err := s.repo.GetById(ctx, userId)
	if err != nil {
		return err
	}

	oldPasswordHash, err := s.hasher.Hash(oldPassword)
	if err != nil {
		return err
	}

	if oldPasswordHash != user.Password {
		return models.ErrOldPassword
	}

	newPasswordHash, err := s.hasher.Hash(newPassword)
	if err != nil {
		return err
	}

	if err := s.repo.UpdateField(ctx, "password", newPasswordHash, userId); err != nil {
		return err
	}

	return nil
}

func (s *UsersService) UpdateInfo(ctx context.Context, userId int, login, firstName, lastName, phoneCode, phoneNumber string) error {
	user, err := s.repo.GetById(ctx, userId)
	if err != nil {
		return err
	}

	if user.Login != login {
		if err := s.repo.UpdateField(ctx, "login", login, userId); err != nil {
			return err
		}
	}

	if err := s.repo.UpdateField(ctx, "first_name", firstName, userId); err != nil {
		return err
	}

	if err := s.repo.UpdateField(ctx, "last_name", lastName, userId); err != nil {
		return err
	}

	if err := s.repo.UpdatePhone(ctx, phoneCode, phoneNumber, userId); err != nil {
		return err
	}

	return nil
}

func (s *UsersService) UpdateAddress(ctx context.Context, userId int, different bool, invoiceAddress models.Address, shippingAddress models.Address) error {
	if different {
		newInvoiceAddress, err := s.repo.CreateAddress(ctx, invoiceAddress)
		if err != nil {
			return err
		}

		if err := s.repo.LinkAddress(ctx, "invoice", userId, newInvoiceAddress.Id); err != nil {
			return err
		}

		newShippingAddress, err := s.repo.CreateAddress(ctx, shippingAddress)
		if err != nil {
			return err
		}

		if err := s.repo.LinkAddress(ctx, "shipping", userId, newShippingAddress.Id); err != nil {
			return err
		}
	} else {
		address, err := s.repo.CreateAddress(ctx, invoiceAddress)
		if err != nil {
			return err
		}

		if err := s.repo.LinkAddress(ctx, "invoice", userId, address.Id); err != nil {
			return err
		}

		if err := s.repo.LinkAddress(ctx, "shipping", userId, address.Id); err != nil {
			return err
		}
	}

	return nil
}
