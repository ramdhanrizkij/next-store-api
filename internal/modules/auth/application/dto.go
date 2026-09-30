package application

import (
	"time"

	userApp "github.com/ramdhanrizkij/next-store-api/internal/modules/user/application"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}

type RegisterResponse struct {
	User    *userApp.UserResponse `json:"user"`
	Message string                `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User        *userApp.UserResponse `json:"user"`
	AccessToken string                `json:"access_token"`
	TokenType   string                `json:"token_type"`
	ExpiresAt   time.Time             `json:"expires_at"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}
