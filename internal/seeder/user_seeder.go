package seeder

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userSeeder struct{}

func NewUserSeeder() Seeder {
	return &userSeeder{}
}

func (s *userSeeder) Name() string {
	return "UserSeeder"
}

func (s *userSeeder) Seed(ctx context.Context, db *sql.DB) error {
	usersToSeed := []struct {
		Name     string
		Email    string
		Password string
		Role     string
	}{
		{
			Name:     "System Admin",
			Email:    "admin@nextstore.com",
			Password: "adminpassword123",
			Role:     "admin",
		},
		{
			Name:     "Demo Customer",
			Email:    "customer@nextstore.com",
			Password: "customerpassword123",
			Role:     "user",
		},
	}

	query := `
		INSERT INTO users (id, name, email, password, role, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (email) DO NOTHING
	`

	for _, u := range usersToSeed {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password for %s: %w", u.Email, err)
		}

		now := time.Now()
		result, err := db.ExecContext(ctx, query,
			uuid.NewString(),
			u.Name,
			u.Email,
			string(hashedPassword),
			u.Role,
			true, // Pre-verified so demo accounts can log in immediately
			now,
			now,
		)
		if err != nil {
			return fmt.Errorf("failed to seed user %s: %w", u.Email, err)
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			log.Printf("   [CREATED] User '%s' (%s, role: %s)", u.Name, u.Email, u.Role)
		} else {
			log.Printf("   [SKIPPED] User '%s' already exists", u.Email)
		}
	}

	return nil
}
