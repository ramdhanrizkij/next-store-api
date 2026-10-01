package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type UserSession struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	UserID    string    `json:"user_id" gorm:"column:user_id"`
	Token     string    `json:"token" gorm:"column:token;uniqueIndex"`
	ExpiresAt time.Time `json:"expires_at" gorm:"column:expires_at"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}

type EmailVerification struct {
	ID         int64      `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	UserID     string     `json:"user_id" gorm:"column:user_id"`
	Token      string     `json:"token" gorm:"column:token;uniqueIndex"`
	ExpiresAt  time.Time  `json:"expires_at" gorm:"column:expires_at"`
	CreatedAt  time.Time  `json:"created_at" gorm:"column:created_at"`
	VerifiedAt *time.Time `json:"verified_at,omitempty" gorm:"column:verified_at"`
}

func (EmailVerification) TableName() string {
	return "email_verifications"
}
