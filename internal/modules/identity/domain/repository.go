package domain

import (
	"context"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindAll(ctx context.Context, limit, offset int) ([]*User, int64, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	SetVerified(ctx context.Context, id string) error
	UpdateLastLogin(ctx context.Context, id string) error
	AssignRole(ctx context.Context, userID, roleID string) error
	FindRoleByName(ctx context.Context, name string) (*Role, error)
}

type AuthRepository interface {
	StoreSession(ctx context.Context, userID, token string, expiresAt time.Time) error
	IsSessionValid(ctx context.Context, userID, token string) (bool, error)
	RevokeSession(ctx context.Context, token string) error

	CreateEmailVerification(ctx context.Context, v *EmailVerification) error
	GetEmailVerificationByToken(ctx context.Context, token string) (*EmailVerification, error)
	MarkEmailVerificationUsed(ctx context.Context, token string) error
}
