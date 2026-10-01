package infrastructure

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ramdhanrizkij/next-store-api/internal/modules/identity/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userPostgresRepository struct {
	db *gorm.DB
}

func NewUserPostgresRepository(db *gorm.DB) domain.UserRepository {
	return &userPostgresRepository{db: db}
}

func (r *userPostgresRepository) Create(ctx context.Context, u *domain.User) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *userPostgresRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var u domain.User
	if err := r.db.WithContext(ctx).Preload("Roles").First(&u, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userPostgresRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	if err := r.db.WithContext(ctx).Preload("Roles").First(&u, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userPostgresRepository) FindAll(ctx context.Context, limit, offset int) ([]*domain.User, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []*domain.User
	if err := r.db.WithContext(ctx).
		Preload("Roles").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userPostgresRepository) Update(ctx context.Context, u *domain.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *userPostgresRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.User{}, "id = ?", id).Error
}

func (r *userPostgresRepository) SetVerified(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&domain.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"email_verified_at": &now,
			"status":            domain.UserStatusActive,
			"updated_at":        now,
		}).Error
}

func (r *userPostgresRepository) UpdateLastLogin(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&domain.User{}).
		Where("id = ?", id).
		Update("last_login_at", &now).Error
}

func (r *userPostgresRepository) AssignRole(ctx context.Context, userID, roleID string) error {
	userRole := domain.UserRole{
		ID:        uuid.NewString(),
		UserID:    userID,
		RoleID:    roleID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&userRole).Error
}

func (r *userPostgresRepository) FindRoleByName(ctx context.Context, name string) (*domain.Role, error) {
	var role domain.Role
	if err := r.db.WithContext(ctx).First(&role, "name = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}
