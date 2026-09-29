package infrastructure

import (
	"context"
	"database/sql"
	"time"

	"github.com/ramdhanrizkij/next-store-api/internal/modules/auth/domain"
)

type authPostgresRepository struct {
	db *sql.DB
}

func NewAuthPostgresRepository(db *sql.DB) domain.AuthRepository {
	return &authPostgresRepository{db: db}
}

func (r *authPostgresRepository) StoreSession(ctx context.Context, userID, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO user_sessions (user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (token) DO UPDATE SET expires_at = EXCLUDED.expires_at
	`
	_, err := r.db.ExecContext(ctx, query, userID, token, expiresAt)
	return err
}

func (r *authPostgresRepository) IsSessionValid(ctx context.Context, userID, token string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM user_sessions 
			WHERE user_id = $1 AND token = $2 AND expires_at > NOW()
		)
	`
	var valid bool
	err := r.db.QueryRowContext(ctx, query, userID, token).Scan(&valid)
	return valid, err
}

func (r *authPostgresRepository) RevokeSession(ctx context.Context, token string) error {
	query := `DELETE FROM user_sessions WHERE token = $1`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}
