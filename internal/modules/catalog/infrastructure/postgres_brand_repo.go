package infrastructure

import (
	"context"
	"errors"

	"github.com/ramdhanrizkij/next-store-api/internal/modules/catalog/domain"
	"gorm.io/gorm"
)

type brandPostgresRepository struct {
	db *gorm.DB
}

func NewBrandPostgresRepository(db *gorm.DB) domain.BrandRepository {
	return &brandPostgresRepository{db: db}
}

// Create implements [domain.BrandRepository].
func (r *brandPostgresRepository) Create(ctx context.Context, brand *domain.Brand) error {
	return r.db.WithContext(ctx).Create(brand).Error
}

// FindByID implements [domain.BrandRepository].
func (r *brandPostgresRepository) FindByID(ctx context.Context, id string) (*domain.Brand, error) {
	var brand domain.Brand
	if err := r.db.WithContext(ctx).First(&brand, "id=?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &brand, nil
}

// FindBySlug implements [domain.BrandRepository].
func (r *brandPostgresRepository) FindBySlug(ctx context.Context, slug string) (*domain.Brand, error) {
	var brand domain.Brand
	if err := r.db.WithContext(ctx).Where("slug=?", slug).First(&brand).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &brand, nil
}

// FindAll implements [domain.BrandRepository].
func (r *brandPostgresRepository) FindAll(ctx context.Context, filter domain.BrandFilter) ([]*domain.Brand, int64, error) {
	query := r.db.WithContext(ctx).Model(&domain.Brand{})

	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR slug ILIKE ?", searchPattern, searchPattern)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 10
	}

	var brands []*domain.Brand
	if err := query.Order("created_at DESC").
		Limit(limit).
		Offset(filter.Offset).
		Find(&brands).Error; err != nil {
		return nil, 0, err
	}

	return brands, total, nil
}

// Update implements domain.BrandRepository.
func (r *brandPostgresRepository) Update(
	ctx context.Context,
	brand *domain.Brand,
) error {
	result := r.db.
		WithContext(ctx).
		Model(&domain.Brand{}).
		Where("id = ?", brand.ID).
		Updates(map[string]interface{}{
			"name":        brand.Name,
			"slug":        brand.Slug,
			"description": brand.Description,
			"logo_url":    brand.LogoUrl,
			"is_active":   brand.IsActive,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Delete implements [domain.BrandRepository].
func (r *brandPostgresRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Where("id=?", id).
		Delete(&domain.Brand{}).
		Error
}
