package application

import (
	"time"

	"github.com/ramdhanrizkij/next-store-api/internal/modules/user/domain"
)

type UserResponse struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Email     string      `json:"email"`
	Role      domain.Role `json:"role"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type UpdateUserRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

func ToUserResponse(user *domain.User) *UserResponse {
	if user == nil {
		return nil
	}
	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserResponses(users []*domain.User) []*UserResponse {
	responses := make([]*UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, ToUserResponse(u))
	}
	return responses
}
