package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ramdhanrizkij/next-store-api/internal/modules/user/domain"
)

type userPostgresRepository struct {
	db *sql.DB
}

func NewUserPostgresRepository(db *sql.DB) domain.UserRepository {
	return &userPostgresRepository{db: db}
}

func (r *userPostgresRepository) Create(ctx context.Context, u *domain.User) error {
	query := `
		INSERT INTO users (id, name, email, password, role, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		u.ID,
		u.Name,
		u.Email,
		u.Password,
		u.Role,
		u.IsVerified,
		u.CreatedAt,
		u.UpdatedAt,
	)
	return err
}

func (r *userPostgresRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, name, email, password, role, is_verified, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var u domain.User
	err := row.Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Password,
		&u.Role,
		&u.IsVerified,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (r *userPostgresRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, name, email, password, role, is_verified, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	row := r.db.QueryRowContext(ctx, query, email)

	var u domain.User
	err := row.Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Password,
		&u.Role,
		&u.IsVerified,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (r *userPostgresRepository) FindAll(ctx context.Context, limit, offset int) ([]*domain.User, int64, error) {
	countQuery := `SELECT COUNT(*) FROM users`
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, name, email, password, role, is_verified, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]*domain.User, 0)
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.Password,
			&u.Role,
			&u.IsVerified,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userPostgresRepository) Update(ctx context.Context, u *domain.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, password = $3, role = $4, is_verified = $5, updated_at = $6
		WHERE id = $7
	`
	_, err := r.db.ExecContext(ctx, query,
		u.Name,
		u.Email,
		u.Password,
		u.Role,
		u.IsVerified,
		u.UpdatedAt,
		u.ID,
	)
	return err
}

func (r *userPostgresRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *userPostgresRepository) SetVerified(ctx context.Context, id string) error {
	query := `UPDATE users SET is_verified = TRUE, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
