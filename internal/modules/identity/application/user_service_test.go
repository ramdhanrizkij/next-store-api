package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/ramdhanrizkij/next-store-api/internal/modules/identity/application"
	"github.com/ramdhanrizkij/next-store-api/internal/modules/identity/domain"
	appErrors "github.com/ramdhanrizkij/next-store-api/internal/shared/errors"
	"github.com/ramdhanrizkij/next-store-api/internal/shared/pagination"
)

type mockUserRepo struct {
	users map[string]*domain.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*domain.User)}
}

func (m *mockUserRepo) Create(ctx context.Context, user *domain.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	u, exists := m.users[id]
	if !exists {
		return nil, nil
	}
	return u, nil
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepo) FindAll(ctx context.Context, limit, offset int) ([]*domain.User, int64, error) {
	var list []*domain.User
	for _, u := range m.users {
		list = append(list, u)
	}
	return list, int64(len(list)), nil
}

func (m *mockUserRepo) Update(ctx context.Context, user *domain.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) Delete(ctx context.Context, id string) error {
	delete(m.users, id)
	return nil
}

func (m *mockUserRepo) SetVerified(ctx context.Context, id string) error {
	u, exists := m.users[id]
	if !exists {
		return appErrors.ErrNotFound
	}
	now := time.Now()
	u.EmailVerifiedAt = &now
	u.Status = domain.UserStatusActive
	return nil
}

func (m *mockUserRepo) UpdateLastLogin(ctx context.Context, id string) error {
	u, exists := m.users[id]
	if !exists {
		return appErrors.ErrNotFound
	}
	now := time.Now()
	u.LastLoginAt = &now
	return nil
}

func (m *mockUserRepo) AssignRole(ctx context.Context, userID, roleID string) error {
	return nil
}

func (m *mockUserRepo) FindRoleByName(ctx context.Context, name string) (*domain.Role, error) {
	return &domain.Role{ID: "r1", Name: name}, nil
}

func TestUserService_GetByID(t *testing.T) {
	repo := newMockUserRepo()
	svc := application.NewUserService(repo)

	user := &domain.User{
		ID:        "u1",
		Name:      "Alice",
		Email:     "alice@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_ = repo.Create(context.Background(), user)

	t.Run("success", func(t *testing.T) {
		res, err := svc.GetByID(context.Background(), "u1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Name != "Alice" {
			t.Errorf("expected Alice, got %s", res.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := svc.GetByID(context.Background(), "unknown")
		if err != appErrors.ErrNotFound {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestUserService_List(t *testing.T) {
	repo := newMockUserRepo()
	svc := application.NewUserService(repo)

	_ = repo.Create(context.Background(), &domain.User{ID: "u1", Name: "Alice"})
	_ = repo.Create(context.Background(), &domain.User{ID: "u2", Name: "Bob"})

	p := pagination.Pagination{Page: 1, Limit: 10}
	users, meta, err := svc.List(context.Background(), p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
	if meta.TotalItems != 2 {
		t.Errorf("expected total 2, got %d", meta.TotalItems)
	}
}
