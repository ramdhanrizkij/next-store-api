package infrastructure

import (
	"context"
	"errors"
	"time"

	"github.com/ramdhanrizkij/next-store-api/internal/modules/identity/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type authPostgresRepository struct {
	db *gorm.DB
}

func NewAuthPostgresRepository(db *gorm.DB) domain.AuthRepository {
	return &authPostgresRepository{db: db}
}

func (r *authPostgresRepository) StoreSession(ctx context.Context, userID, token string, expiresAt time.Time) error {
	session := domain.UserSession{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "token"}},
		DoUpdates: clause.AssignmentColumns([]string{"expires_at"}),
	}).Create(&session).Error
}

func (r *authPostgresRepository) IsSessionValid(ctx context.Context, userID, token string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.UserSession{}).
		Where("user_id = ? AND token = ? AND expires_at > ?", userID, token, time.Now()).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *authPostgresRepository) RevokeSession(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Where("token = ?", token).Delete(&domain.UserSession{}).Error
}

func (r *authPostgresRepository) CreateEmailVerification(ctx context.Context, v *domain.EmailVerification) error {
	return r.db.WithContext(ctx).Create(v).Error
}

func (r *authPostgresRepository) GetEmailVerificationByToken(ctx context.Context, token string) (*domain.EmailVerification, error) {
	var v domain.EmailVerification
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&v).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &v, nil
}

func (r *authPostgresRepository) MarkEmailVerificationUsed(ctx context.Context, token string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&domain.EmailVerification{}).
		Where("token = ?", token).
		Update("verified_at", &now).Error
}
