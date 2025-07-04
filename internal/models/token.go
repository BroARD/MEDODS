package models

import "time"

type RefreshToken struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	RefreshTokenHash []byte    `json:"refresh_token_hash"`
	UserAgent        string    `json:"user_agent"`
	IP               string    `json:"ip"`
	CreatedAt        time.Time `json:"created_at"`
	ExpiresAt        time.Time `json:"expires_at"`
}
