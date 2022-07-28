package models

import "time"

type Image struct {
	Id        int       `json:"id,omitempty" db:"id"`
	Filename  string    `json:"filename" db:"filename"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
