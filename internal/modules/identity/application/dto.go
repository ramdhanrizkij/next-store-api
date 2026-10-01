package application

import (
	"time"

	"github.com/ramdhanrizkij/next-store-api/internal/modules/identity/domain"
)

type RegisterRequest struct {
	Name     string  `json:"name" binding:"required,min=2,max=100"`
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6,max=100"`
	Phone    *string `json:"phone,omitempty"`
}

type RegisterResponse struct {
	User    *UserResponse `json:"user"`
	Message string        `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User        *UserResponse `json:"user"`
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
	ExpiresAt   time.Time     `json:"expires_at"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type UserResponse struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Email           string            `json:"email"`
	Phone           *string           `json:"phone,omitempty"`
	Status          domain.UserStatus `json:"status"`
	Role            string            `json:"role"`
	Roles           []string          `json:"roles,omitempty"`
	IsVerified      bool              `json:"is_verified"`
	EmailVerifiedAt *time.Time        `json:"email_verified_at,omitempty"`
	LastLoginAt     *time.Time        `json:"last_login_at,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

type UpdateUserRequest struct {
	Name  string  `json:"name" binding:"required,min=2,max=100"`
	Phone *string `json:"phone,omitempty"`
}

func ToUserResponse(user *domain.User) *UserResponse {
	if user == nil {
		return nil
	}

	roleNames := make([]string, 0, len(user.Roles))
	for _, r := range user.Roles {
		roleNames = append(roleNames, r.Name)
	}

	return &UserResponse{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		Phone:           user.Phone,
		Status:          user.Status,
		Role:            user.PrimaryRole(),
		Roles:           roleNames,
		IsVerified:      user.IsVerified(),
		EmailVerifiedAt: user.EmailVerifiedAt,
		LastLoginAt:     user.LastLoginAt,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}

func ToUserResponses(users []*domain.User) []*UserResponse {
	responses := make([]*UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, ToUserResponse(u))
	}
	return responses
}
