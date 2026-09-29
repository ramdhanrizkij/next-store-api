package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/ramdhanrizkij/next-store-api/internal/modules/auth/application"
	authDomain "github.com/ramdhanrizkij/next-store-api/internal/modules/auth/domain"
	userDomain "github.com/ramdhanrizkij/next-store-api/internal/modules/user/domain"
)

type mockUserRepoForAuth struct {
	users map[string]*userDomain.User
}

func newMockUserRepoForAuth() *mockUserRepoForAuth {
	return &mockUserRepoForAuth{users: make(map[string]*userDomain.User)}
}

func (m *mockUserRepoForAuth) Create(ctx context.Context, user *userDomain.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepoForAuth) FindByID(ctx context.Context, id string) (*userDomain.User, error) {
	u, exists := m.users[id]
	if !exists {
		return nil, nil
	}
	return u, nil
}

func (m *mockUserRepoForAuth) FindByEmail(ctx context.Context, email string) (*userDomain.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepoForAuth) FindAll(ctx context.Context, limit, offset int) ([]*userDomain.User, int64, error) {
	return nil, 0, nil
}

func (m *mockUserRepoForAuth) Update(ctx context.Context, user *userDomain.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepoForAuth) Delete(ctx context.Context, id string) error {
	delete(m.users, id)
	return nil
}

type mockAuthRepo struct{}

func (m *mockAuthRepo) StoreSession(ctx context.Context, userID, token string, expiresAt time.Time) error {
	return nil
}
func (m *mockAuthRepo) IsSessionValid(ctx context.Context, userID, token string) (bool, error) {
	return true, nil
}
func (m *mockAuthRepo) RevokeSession(ctx context.Context, token string) error {
	return nil
}

var _ authDomain.AuthRepository = (*mockAuthRepo)(nil)

func TestAuthService_RegisterAndLogin(t *testing.T) {
	userRepo := newMockUserRepoForAuth()
	authRepo := &mockAuthRepo{}
	svc := application.NewAuthService(userRepo, authRepo, "testsecretkey", 24)

	ctx := context.Background()

	t.Run("Register success", func(t *testing.T) {
		req := &application.RegisterRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "securepassword",
		}

		res, err := svc.Register(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error registering: %v", err)
		}
		if res.AccessToken == "" {
			t.Error("expected non-empty access token")
		}
		if res.User.Email != "test@example.com" {
			t.Errorf("expected email test@example.com, got %s", res.User.Email)
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

	t.Run("Login success", func(t *testing.T) {
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
