package seeder

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	identityDomain "github.com/ramdhanrizkij/next-store-api/internal/modules/identity/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userSeeder struct{}

func NewUserSeeder() Seeder {
	return &userSeeder{}
}

func (s *userSeeder) Name() string {
	return "UserSeeder"
}

func (s *userSeeder) Seed(ctx context.Context, db *gorm.DB) error {
	now := time.Now()

	// 1. Seed Roles
	rolesToSeed := []identityDomain.Role{
		{
			ID:        uuid.NewString(),
			Name:      identityDomain.RoleAdmin,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        uuid.NewString(),
			Name:      identityDomain.RoleUser,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	roleMap := make(map[string]string) // name -> id
	for _, r := range rolesToSeed {
		role := r
		err := db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).Create(&role).Error
		if err != nil {
			return fmt.Errorf("failed to seed role %s: %w", r.Name, err)
		}

		// Retrieve the role to get its real ID in case it already existed
		var existingRole identityDomain.Role
		if err := db.WithContext(ctx).First(&existingRole, "name = ?", r.Name).Error; err == nil {
			roleMap[existingRole.Name] = existingRole.ID
		}
	}

	// 2. Seed Users
	usersToSeed := []struct {
		Name     string
		Email    string
		Password string
		RoleName string
	}{
		{
			Name:     "System Admin",
			Email:    "admin@nextstore.com",
			Password: "adminpassword123",
			RoleName: identityDomain.RoleAdmin,
		},
		{
			Name:     "Demo Customer",
			Email:    "customer@nextstore.com",
			Password: "customerpassword123",
			RoleName: identityDomain.RoleUser,
		},
	}

	for _, u := range usersToSeed {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password for %s: %w", u.Email, err)
		}

		userID := uuid.NewString()
		user := identityDomain.User{
			ID:              userID,
			Name:            u.Name,
			Email:           u.Email,
			Password:        string(hashedPassword),
			Status:          identityDomain.UserStatusActive,
			EmailVerifiedAt: &now,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		result := db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "email"}},
			DoNothing: true,
		}).Create(&user)
		if result.Error != nil {
			return fmt.Errorf("failed to seed user %s: %w", u.Email, result.Error)
		}

		// Retrieve actual user ID
		var existingUser identityDomain.User
		if err := db.WithContext(ctx).First(&existingUser, "email = ?", u.Email).Error; err == nil {
			userID = existingUser.ID
		}

		// Link role in user_roles
		if roleID, ok := roleMap[u.RoleName]; ok {
			userRole := identityDomain.UserRole{
				ID:        uuid.NewString(),
				UserID:    userID,
				RoleID:    roleID,
				CreatedAt: now,
				UpdatedAt: now,
			}
			_ = db.WithContext(ctx).Clauses(clause.OnConflict{
				DoNothing: true,
			}).Create(&userRole)
		}

		if result.RowsAffected > 0 {
			log.Printf("   [CREATED] User '%s' (%s, role: %s)", u.Name, u.Email, u.RoleName)
		} else {
			log.Printf("   [SKIPPED] User '%s' already exists", u.Email)
		}
	}

	return nil
}
