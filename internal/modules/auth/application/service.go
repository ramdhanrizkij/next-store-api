package application

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/ramdhanrizkij/next-store-api/internal/modules/auth/domain"
	userApp "github.com/ramdhanrizkij/next-store-api/internal/modules/user/application"
	userDomain "github.com/ramdhanrizkij/next-store-api/internal/modules/user/domain"
	appErrors "github.com/ramdhanrizkij/next-store-api/internal/shared/errors"
)

type AuthService interface {
	Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
}

type authService struct {
	userRepo       userDomain.UserRepository
	authRepo       domain.AuthRepository
	jwtSecret      string
	jwtExpiryHours int
}

func NewAuthService(
	userRepo userDomain.UserRepository,
	authRepo domain.AuthRepository,
	jwtSecret string,
	jwtExpiryHours int,
) AuthService {
	return &authService{
		userRepo:       userRepo,
		authRepo:       authRepo,
		jwtSecret:      jwtSecret,
		jwtExpiryHours: jwtExpiryHours,
	}
}

func (s *authService) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
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
		ID:        uuid.NewString(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Role:      userDomain.RoleUser,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	token, expiresAt, err := s.generateToken(newUser)
	if err != nil {
		return nil, err
	}

	if s.authRepo != nil {
		_ = s.authRepo.StoreSession(ctx, newUser.ID, token, expiresAt)
	}

	return &AuthResponse{
		User:        userApp.ToUserResponse(newUser),
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt,
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
