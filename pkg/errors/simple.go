package errors

import "errors"

var (
	// HTTP errors
	ErrInvalidBody     = errors.New("invalid request body")
	ErrInvalidFormBody = errors.New("invalid request form body")
	ErrInvalidParam    = errors.New("invalid request param")
	ErrInvalidCookie   = errors.New("invalid request cookie")

	// Auth errors
	ErrEmptyAuthHeader   = errors.New("empty auth header")
	ErrInvalidAuthHeader = errors.New("invalid auth header")

	// Repo errors
	ErrUserNotFound     = errors.New("user not found")
	ErrUserNotCompleted = errors.New("user must verify email")
	ErrAddressNotFound  = errors.New("address not found")
	ErrDeliveryNotFound = errors.New("delivery not found")
	ErrOldPassword      = errors.New("wrong old password")
	ErrFileExtension    = errors.New("wrong file extension")
	ErrViolatesKey      = errors.New("violates foreign key constraint")

	// Sort errors
	ErrSortOptions = errors.New("wrong sort options")

	// SMTP
	ErrEmailSend = errors.New("cannot send email")
)
