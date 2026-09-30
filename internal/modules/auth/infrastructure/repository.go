package infrastructure

import (
	"context"
	"database/sql"
	"errors"
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

func (r *authPostgresRepository) CreateEmailVerification(ctx context.Context, v *domain.EmailVerification) error {
	query := `
		INSERT INTO email_verifications (user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		v.UserID,
		v.Token,
		v.ExpiresAt,
		v.CreatedAt,
	).Scan(&v.ID)
}

func (r *authPostgresRepository) GetEmailVerificationByToken(ctx context.Context, token string) (*domain.EmailVerification, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at, verified_at
		FROM email_verifications
		WHERE token = $1
	`
	row := r.db.QueryRowContext(ctx, query, token)

	var v domain.EmailVerification
	var verifiedAt sql.NullTime
	err := row.Scan(
		&v.ID,
		&v.UserID,
		&v.Token,
		&v.ExpiresAt,
		&v.CreatedAt,
		&verifiedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if verifiedAt.Valid {
		v.VerifiedAt = &verifiedAt.Time
	}

	return &v, nil
}

func (r *authPostgresRepository) MarkEmailVerificationUsed(ctx context.Context, token string) error {
	query := `
		UPDATE email_verifications
		SET verified_at = NOW()
		WHERE token = $1
	`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}
