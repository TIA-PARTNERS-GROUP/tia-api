package ports

import "time"

type LoginInput struct {
	LoginEmail string `json:"login_email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User      UserResponse `json:"user"`
	Token     string       `json:"token"`
	SessionID uint         `json:"session_id"`
	ExpiresAt time.Time    `json:"expires_at"`
	TokenType string       `json:"token_type"`
}
