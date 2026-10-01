package domain

import "context"

type BrandFilter struct {
	Search   string
	IsActive *bool
	Limit    int
	Offset   int
}

type BrandRepository interface {
	Create(ctx context.Context, brand *Brand) error
	FindByID(ctx context.Context, id string) (*Brand, error)
	FindBySlug(ctx context.Context, slug string) (*Brand, error)
	FindAll(ctx context.Context, filter BrandFilter) ([]*Brand, int64, error)
	Update(ctx context.Context, brand *Brand) error
	Delete(ctx context.Context, id string) error
}
