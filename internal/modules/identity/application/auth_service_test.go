package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/hibiken/asynq"
	"github.com/ramdhanrizkij/next-store-api/internal/modules/identity/application"
	"github.com/ramdhanrizkij/next-store-api/internal/modules/identity/domain"
	appErrors "github.com/ramdhanrizkij/next-store-api/internal/shared/errors"
	"github.com/ramdhanrizkij/next-store-api/internal/worker"
)

type mockUserRepoForAuth struct {
	users map[string]*domain.User
}

func newMockUserRepoForAuth() *mockUserRepoForAuth {
	return &mockUserRepoForAuth{users: make(map[string]*domain.User)}
}

func (m *mockUserRepoForAuth) Create(ctx context.Context, user *domain.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepoForAuth) FindByID(ctx context.Context, id string) (*domain.User, error) {
	u, exists := m.users[id]
	if !exists {
		return nil, nil
	}
	return u, nil
}

func (m *mockUserRepoForAuth) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepoForAuth) FindAll(ctx context.Context, limit, offset int) ([]*domain.User, int64, error) {
	return nil, 0, nil
}

func (m *mockUserRepoForAuth) Update(ctx context.Context, user *domain.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepoForAuth) Delete(ctx context.Context, id string) error {
	delete(m.users, id)
	return nil
}

func (m *mockUserRepoForAuth) SetVerified(ctx context.Context, id string) error {
	u, exists := m.users[id]
	if !exists {
		return appErrors.ErrNotFound
	}
	now := time.Now()
	u.EmailVerifiedAt = &now
	u.Status = domain.UserStatusActive
	return nil
}

func (m *mockUserRepoForAuth) UpdateLastLogin(ctx context.Context, id string) error {
	u, exists := m.users[id]
	if !exists {
		return appErrors.ErrNotFound
	}
	now := time.Now()
	u.LastLoginAt = &now
	return nil
}

func (m *mockUserRepoForAuth) AssignRole(ctx context.Context, userID, roleID string) error {
	return nil
}

func (m *mockUserRepoForAuth) FindRoleByName(ctx context.Context, name string) (*domain.Role, error) {
	return &domain.Role{ID: "r1", Name: name}, nil
}

type mockAuthRepo struct {
	verifications map[string]*domain.EmailVerification
}

func newMockAuthRepo() *mockAuthRepo {
	return &mockAuthRepo{
		verifications: make(map[string]*domain.EmailVerification),
	}
}

func (m *mockAuthRepo) StoreSession(ctx context.Context, userID, token string, expiresAt time.Time) error {
	return nil
}

func (m *mockAuthRepo) IsSessionValid(ctx context.Context, userID, token string) (bool, error) {
	return true, nil
}

func (m *mockAuthRepo) RevokeSession(ctx context.Context, token string) error {
	return nil
}

func (m *mockAuthRepo) CreateEmailVerification(ctx context.Context, v *domain.EmailVerification) error {
	m.verifications[v.Token] = v
	return nil
}

func (m *mockAuthRepo) GetEmailVerificationByToken(ctx context.Context, token string) (*domain.EmailVerification, error) {
	v, exists := m.verifications[token]
	if !exists {
		return nil, nil
	}
	return v, nil
}

func (m *mockAuthRepo) MarkEmailVerificationUsed(ctx context.Context, token string) error {
	v, exists := m.verifications[token]
	if exists {
		now := time.Now()
		v.VerifiedAt = &now
	}
	return nil
}

var _ domain.AuthRepository = (*mockAuthRepo)(nil)

type mockTaskDistributor struct {
	tasks []*worker.SendEmailVerificationPayload
}

func (m *mockTaskDistributor) DistributeTaskSendEmailVerification(
	ctx context.Context,
	payload *worker.SendEmailVerificationPayload,
	opts ...asynq.Option,
) error {
	m.tasks = append(m.tasks, payload)
	return nil
}

func TestAuthService_RegisterVerificationAndLogin(t *testing.T) {
	userRepo := newMockUserRepoForAuth()
	authRepo := newMockAuthRepo()
	taskDist := &mockTaskDistributor{}
	svc := application.NewAuthService(
		userRepo,
		authRepo,
		taskDist,
		"http://localhost:8080",
		"testsecretkey",
		24,
	)

	ctx := context.Background()

	var verificationToken string

	t.Run("Register success - creates unverified user and enqueues task", func(t *testing.T) {
		req := &application.RegisterRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "securepassword",
		}

		res, err := svc.Register(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error registering: %v", err)
		}
		if res.User == nil {
			t.Fatal("expected non-nil user response")
		}
		if res.User.Email != "test@example.com" {
			t.Errorf("expected email test@example.com, got %s", res.User.Email)
		}
		if res.User.IsVerified {
			t.Error("expected user to not be verified yet")
		}

		if len(taskDist.tasks) != 1 {
			t.Fatalf("expected 1 task distributed, got %d", len(taskDist.tasks))
		}
		verificationToken = taskDist.tasks[0].Token
		if verificationToken == "" {
			t.Fatal("expected non-empty verification token in task payload")
		}
	})

	t.Run("Register duplicate email fails", func(t *testing.T) {
		req := &application.RegisterRequest{
			Name:     "Duplicate",
			Email:    "test@example.com",
			Password: "securepassword",
		}

		_, err := svc.Register(ctx, req)
		if err == nil {
			t.Fatal("expected error for duplicate email, got nil")
		}
	})

	t.Run("Login unverified account fails", func(t *testing.T) {
		req := &application.LoginRequest{
			Email:    "test@example.com",
			Password: "securepassword",
		}

		_, err := svc.Login(ctx, req)
		if err == nil {
			t.Fatal("expected error for unverified account login, got nil")
		}
	})

	t.Run("VerifyEmail with invalid token fails", func(t *testing.T) {
		err := svc.VerifyEmail(ctx, &application.VerifyEmailRequest{
			Token: "invalid-token",
		})
		if err == nil {
			t.Fatal("expected error for invalid token, got nil")
		}
	})

	t.Run("VerifyEmail with valid token succeeds", func(t *testing.T) {
		err := svc.VerifyEmail(ctx, &application.VerifyEmailRequest{
			Token: verificationToken,
		})
		if err != nil {
			t.Fatalf("unexpected error verifying email: %v", err)
		}
	})

	t.Run("VerifyEmail already verified token fails", func(t *testing.T) {
		err := svc.VerifyEmail(ctx, &application.VerifyEmailRequest{
			Token: verificationToken,
		})
		if err == nil {
			t.Fatal("expected error for already verified token, got nil")
		}
	})

	t.Run("Login success after verification", func(t *testing.T) {
		req := &application.LoginRequest{
			Email:    "test@example.com",
			Password: "securepassword",
		}

		res, err := svc.Login(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error logging in: %v", err)
		}
		if res.AccessToken == "" {
			t.Error("expected non-empty access token")
		}
	})

	t.Run("Login invalid password fails", func(t *testing.T) {
		req := &application.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		_, err := svc.Login(ctx, req)
		if err == nil {
			t.Fatal("expected error for wrong password, got nil")
		}
	})
}
