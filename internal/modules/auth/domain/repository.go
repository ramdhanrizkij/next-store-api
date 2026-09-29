package domain

import (
	"context"
	"time"
)

type AuthRepository interface {
	StoreSession(ctx context.Context, userID, token string, expiresAt time.Time) error
	IsSessionValid(ctx context.Context, userID, token string) (bool, error)
	RevokeSession(ctx context.Context, token string) error

	CreateEmailVerification(ctx context.Context, v *EmailVerification) error
	GetEmailVerificationByToken(ctx context.Context, token string) (*EmailVerification, error)
	MarkEmailVerificationUsed(ctx context.Context, token string) error
}
