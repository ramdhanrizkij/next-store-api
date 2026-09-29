package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/ramdhanrizkij/next-store-api/internal/modules/auth/domain"
	userApp "github.com/ramdhanrizkij/next-store-api/internal/modules/user/application"
	userDomain "github.com/ramdhanrizkij/next-store-api/internal/modules/user/domain"
	appErrors "github.com/ramdhanrizkij/next-store-api/internal/shared/errors"
	"github.com/ramdhanrizkij/next-store-api/internal/worker"
)

type AuthService interface {
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
	VerifyEmail(ctx context.Context, req *VerifyEmailRequest) error
}

type authService struct {
	userRepo        userDomain.UserRepository
	authRepo        domain.AuthRepository
	taskDistributor worker.TaskDistributor
	appURL          string
	jwtSecret       string
	jwtExpiryHours  int
}

func NewAuthService(
	userRepo userDomain.UserRepository,
	authRepo domain.AuthRepository,
	taskDistributor worker.TaskDistributor,
	appURL string,
	jwtSecret string,
	jwtExpiryHours int,
) AuthService {
	return &authService{
		userRepo:        userRepo,
		authRepo:        authRepo,
		taskDistributor: taskDistributor,
		appURL:          appURL,
		jwtSecret:       jwtSecret,
		jwtExpiryHours:  jwtExpiryHours,
	}
}

func (s *authService) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, appErrors.New(http.StatusConflict, "EMAIL_ALREADY_EXISTS", "A user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, appErrors.Wrap(err, http.StatusInternalServerError, "PASSWORD_HASH_ERROR", "Failed to secure password")
	}

	now := time.Now()
	newUser := &userDomain.User{
		ID:         uuid.NewString(),
		Name:       req.Name,
		Email:      req.Email,
		Password:   string(hashedPassword),
		Role:       userDomain.RoleUser,
		IsVerified: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	// Generate a secure random token for email verification
	verificationToken := generateSecureToken(32)
	expiresAt := now.Add(24 * time.Hour)

	if s.authRepo != nil {
		verification := &domain.EmailVerification{
			UserID:    newUser.ID,
			Token:     verificationToken,
			ExpiresAt: expiresAt,
			CreatedAt: now,
		}
		if err := s.authRepo.CreateEmailVerification(ctx, verification); err != nil {
			return nil, appErrors.Wrap(err, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to record email verification")
		}
	}

	// Dispatch asynchronous email verification task via Redis worker
	if s.taskDistributor != nil {
		verificationURL := fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", s.appURL, verificationToken)
		_ = s.taskDistributor.DistributeTaskSendEmailVerification(ctx, &worker.SendEmailVerificationPayload{
			UserID:          newUser.ID,
			Name:            newUser.Name,
			Email:           newUser.Email,
			Token:           verificationToken,
			VerificationURL: verificationURL,
		})
	}

	return &RegisterResponse{
		User:    userApp.ToUserResponse(newUser),
		Message: "Registration successful. Please check your email to verify and activate your account.",
	}, nil
}

func (s *authService) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, appErrors.New(http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, appErrors.New(http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
	}

	// Disallow login for unverified accounts
	if !user.IsVerified {
		return nil, appErrors.New(http.StatusForbidden, "EMAIL_NOT_VERIFIED", "Your email address is not verified yet. Please check your inbox to activate your account.")
	}

	token, expiresAt, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	if s.authRepo != nil {
		_ = s.authRepo.StoreSession(ctx, user.ID, token, expiresAt)
	}

	return &AuthResponse{
		User:        userApp.ToUserResponse(user),
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt,
	}, nil
}

func (s *authService) VerifyEmail(ctx context.Context, req *VerifyEmailRequest) error {
	if s.authRepo == nil {
		return errors.New("auth repository not configured")
	}

	record, err := s.authRepo.GetEmailVerificationByToken(ctx, req.Token)
	if err != nil {
		return err
	}
	if record == nil {
		return appErrors.New(http.StatusBadRequest, "INVALID_TOKEN", "Verification token is invalid or does not exist")
	}

	if record.VerifiedAt != nil {
		return appErrors.New(http.StatusBadRequest, "ALREADY_VERIFIED", "This email account has already been verified")
	}

	if time.Now().After(record.ExpiresAt) {
		return appErrors.New(http.StatusBadRequest, "TOKEN_EXPIRED", "Verification token has expired. Please request a new verification email.")
	}

	// Mark user as verified
	if err := s.userRepo.SetVerified(ctx, record.UserID); err != nil {
		return appErrors.Wrap(err, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to activate user account")
	}

	// Mark token as used
	if err := s.authRepo.MarkEmailVerificationUsed(ctx, req.Token); err != nil {
		return appErrors.Wrap(err, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to mark token as used")
	}

	return nil
}

func (s *authService) generateToken(user *userDomain.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(s.jwtExpiryHours) * time.Hour)
	claims := domain.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", time.Time{}, errors.New("failed to generate access token")
	}

	return tokenString, expiresAt, nil
}

func generateSecureToken(length int) string {
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
