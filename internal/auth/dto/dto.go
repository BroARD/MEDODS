package dto

type TokensPair struct {
	AccessToken string
	RefreshToken string
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

