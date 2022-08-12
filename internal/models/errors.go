package models

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyAuthHeader   = errors.New("empty auth header")
	ErrInvalidAuthHeader = errors.New("invalid auth header")
	ErrUserNotFound      = errors.New("user not found")
	ErrAddressNotFound   = errors.New("address not found")
	ErrDeliveryNotFound  = errors.New("delivery not found")
	ErrOldPassword       = errors.New("wrong old password")
	ErrFileExtension     = errors.New("wrong file extension")
	ErrViolatesKey       = errors.New("violates foreign key constraint")
)

type ErrUniqueValue struct {
	Field string
}

func (e ErrUniqueValue) Error() string {
	return fmt.Sprintf("%s already exist", e.Field)
}

func NewErrUniqueValue(field string) ErrUniqueValue {
	return ErrUniqueValue{
		Field: field,
	}
}
