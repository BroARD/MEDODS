package dto

type TokensPair struct {
	AccessToken string
	RefreshToken string
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type contextKey string

const UserIDKey = contextKey("user_id")

