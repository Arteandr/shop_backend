package models

type Phone struct {
	Code   *string `json:"phoneCode" db:"code"`
	Number *string `json:"phoneNumber" db:"number"`
}
