package application

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ramdhanrizkij/next-store-api/internal/modules/catalog/domain"
	appErrors "github.com/ramdhanrizkij/next-store-api/internal/shared/errors"
	"github.com/ramdhanrizkij/next-store-api/internal/shared/pagination"
	"github.com/ramdhanrizkij/next-store-api/internal/shared/storage"
)

type BrandService interface {
	Create(ctx context.Context, req *CreateBrandRequest) (*BrandResponse, error)
	GetByID(ctx context.Context, id string) (*BrandResponse, error)
	GetBySlug(ctx context.Context, slug string) (*BrandResponse, error)
	List(ctx context.Context, search string, isActive *bool, p pagination.Pagination) ([]*BrandResponse, *pagination.Meta, error)
	Update(ctx context.Context, id string, req *UpdateBrandRequest) (*BrandResponse, error)
	Delete(ctx context.Context, id string) error
	UploadLogo(ctx context.Context, fileHeader *multipart.FileHeader) (string, error)
}

type brandService struct {
	brandRepo domain.BrandRepository
	storage   storage.StorageService
}

func NewBrandService(brandRepo domain.BrandRepository, storage storage.StorageService) BrandService {
	return &brandService{
		brandRepo: brandRepo,
		storage:   storage,
	}
}

func (s *brandService) Create(ctx context.Context, req *CreateBrandRequest) (*BrandResponse, error) {
	slug := req.Slug
	if slug == "" {
		slug = generateSlug(req.Name)
	} else {
		slug = generateSlug(slug)
	}

	existing, err := s.brandRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, appErrors.New(http.StatusConflict, "SLUG_ALREADY_EXISTS", fmt.Sprintf("Brand with slug '%s' already exists", slug))
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	now := time.Now()
	brand := &domain.Brand{
		ID:          uuid.New(),
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		LogoUrl:     req.LogoUrl,
		IsActive:    isActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.brandRepo.Create(ctx, brand); err != nil {
		return nil, err
	}

	return ToBrandResponse(brand), nil
}

func (s *brandService) GetByID(ctx context.Context, id string) (*BrandResponse, error) {
	brand, err := s.brandRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, appErrors.ErrNotFound
	}
	return ToBrandResponse(brand), nil
}

func (s *brandService) GetBySlug(ctx context.Context, slug string) (*BrandResponse, error) {
	brand, err := s.brandRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, appErrors.ErrNotFound
	}
	return ToBrandResponse(brand), nil
}

func (s *brandService) List(
	ctx context.Context,
	search string,
	isActive *bool,
	p pagination.Pagination,
) ([]*BrandResponse, *pagination.Meta, error) {
	filter := domain.BrandFilter{
		Search:   search,
		IsActive: isActive,
		Limit:    p.Limit,
		Offset:   p.Offset(),
	}

	brands, total, err := s.brandRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	meta := pagination.BuildMeta(total, p.Page, p.Limit)
	return ToBrandResponses(brands), meta, nil
}

func (s *brandService) Update(ctx context.Context, id string, req *UpdateBrandRequest) (*BrandResponse, error) {
	brand, err := s.brandRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, appErrors.ErrNotFound
	}

	slug := req.Slug
	if slug == "" {
		slug = generateSlug(req.Name)
	} else {
		slug = generateSlug(slug)
	}

	if slug != brand.Slug {
		existing, err := s.brandRepo.FindBySlug(ctx, slug)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != brand.ID {
			return nil, appErrors.New(http.StatusConflict, "SLUG_ALREADY_EXISTS", fmt.Sprintf("Brand with slug '%s' already exists", slug))
		}
	}

	brand.Name = req.Name
	brand.Slug = slug
	if req.Description != nil {
		brand.Description = req.Description
	}
	if req.LogoUrl != nil {
		brand.LogoUrl = req.LogoUrl
	}
	if req.IsActive != nil {
		brand.IsActive = *req.IsActive
	}
	brand.UpdatedAt = time.Now()

	if err := s.brandRepo.Update(ctx, brand); err != nil {
		return nil, err
	}

	return ToBrandResponse(brand), nil
}

func (s *brandService) Delete(ctx context.Context, id string) error {
	brand, err := s.brandRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if brand == nil {
		return appErrors.ErrNotFound
	}

	return s.brandRepo.Delete(ctx, id)
}

func (s *brandService) UploadLogo(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	if s.storage == nil {
		return "", appErrors.New(http.StatusInternalServerError, "STORAGE_ERROR", "Storage service is not configured")
	}

	// Validate file size (max 5MB)
	const maxFileSize = 5 * 1024 * 1024
	if fileHeader.Size > maxFileSize {
		return "", appErrors.New(http.StatusBadRequest, "FILE_TOO_LARGE", "File size exceeds 5MB limit")
	}

	// Validate content type / extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedExts := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".webp": "image/webp",
		".svg":  "image/svg+xml",
	}

	contentType, ok := allowedExts[ext]
	if !ok {
		return "", appErrors.New(http.StatusBadRequest, "INVALID_FILE_TYPE", "Only JPG, PNG, WEBP, and SVG images are allowed")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", appErrors.Wrap(err, http.StatusInternalServerError, "FILE_OPEN_ERROR", "Failed to open uploaded file")
	}
	defer file.Close()

	url, err := s.storage.UploadFile(ctx, "brands", fileHeader.Filename, file, fileHeader.Size, contentType)
	if err != nil {
		return "", appErrors.Wrap(err, http.StatusInternalServerError, "UPLOAD_ERROR", "Failed to upload file to storage")
	}

	return url, nil
}

func generateSlug(input string) string {
	slug := strings.ToLower(input)
	// Replace non-alphanumeric characters with hyphens
	reg := regexp.MustCompile("[^a-z0-9]+")
	slug = reg.ReplaceAllString(slug, "-")
	return strings.Trim(slug, "-")
}
