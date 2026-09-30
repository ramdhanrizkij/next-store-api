package seeder

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	userDomain "github.com/ramdhanrizkij/next-store-api/internal/modules/user/domain"
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
	usersToSeed := []struct {
		Name     string
		Email    string
		Password string
		Role     userDomain.Role
	}{
		{
			Name:     "System Admin",
			Email:    "admin@nextstore.com",
			Password: "adminpassword123",
			Role:     userDomain.RoleAdmin,
		},
		{
			Name:     "Demo Customer",
			Email:    "customer@nextstore.com",
			Password: "customerpassword123",
			Role:     userDomain.RoleUser,
		},
	}

	for _, u := range usersToSeed {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password for %s: %w", u.Email, err)
		}

		now := time.Now()
		user := userDomain.User{
			ID:         uuid.NewString(),
			Name:       u.Name,
			Email:      u.Email,
			Password:   string(hashedPassword),
			Role:       u.Role,
			IsVerified: true, // Pre-verified so demo accounts can log in immediately
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		result := db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "email"}},
			DoNothing: true,
		}).Create(&user)
		if result.Error != nil {
			return fmt.Errorf("failed to seed user %s: %w", u.Email, result.Error)
		}

		if result.RowsAffected > 0 {
			log.Printf("   [CREATED] User '%s' (%s, role: %s)", u.Name, u.Email, u.Role)
		} else {
			log.Printf("   [SKIPPED] User '%s' already exists", u.Email)
		}
	}

	return nil
}
