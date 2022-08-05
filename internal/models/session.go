package models

import "time"

type Session struct {
	RefreshToken string    `json:"refreshToken" bd:"refresh_token"`
	ExpiresAt    time.Time `json:"expiresAt" bd:"expires_at"`
}
