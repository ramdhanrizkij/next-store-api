package application

import (
	"time"

	"github.com/google/uuid"
	"github.com/ramdhanrizkij/next-store-api/internal/modules/catalog/domain"
)

type CreateBrandRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=150"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
	LogoUrl     *string `json:"logo_url"`
	IsActive    *bool   `json:"is_active"`
}

type UpdateBrandRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=150"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
	LogoUrl     *string `json:"logo_url"`
	IsActive    *bool   `json:"is_active"`
}

type BrandResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description,omitempty"`
	LogoUrl     *string   `json:"logo_url,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToBrandResponse(brand *domain.Brand) *BrandResponse {
	if brand == nil {
		return nil
	}
	return &BrandResponse{
		ID:          brand.ID,
		Name:        brand.Name,
		Slug:        brand.Slug,
		Description: brand.Description,
		LogoUrl:     brand.LogoUrl,
		IsActive:    brand.IsActive,
		CreatedAt:   brand.CreatedAt,
		UpdatedAt:   brand.UpdatedAt,
	}
}

func ToBrandResponses(brands []*domain.Brand) []*BrandResponse {
	res := make([]*BrandResponse, 0, len(brands))
	for _, b := range brands {
		res = append(res, ToBrandResponse(b))
	}
	return res
}
