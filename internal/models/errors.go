package models

import "errors"

var (
	ErrEmptyAuthHeader   = errors.New("empty auth header")
	ErrInvalidAuthHeader = errors.New("invalid auth header")
	ErrUserNotFound      = errors.New("user not found")
)
