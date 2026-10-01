package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ramdhanrizkij/next-store-api/internal/modules/catalog/application"
	"github.com/ramdhanrizkij/next-store-api/internal/shared/pagination"
	"github.com/ramdhanrizkij/next-store-api/internal/shared/response"
	"github.com/ramdhanrizkij/next-store-api/internal/shared/validator"
)

type BrandHandler struct {
	brandService application.BrandService
}

func NewBrandHandler(brandService application.BrandService) *BrandHandler {
	return &BrandHandler{
		brandService: brandService,
	}
}

// Create handles creating a new brand
func (h *BrandHandler) Create(c *gin.Context) {
	var req application.CreateBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", validator.FormatValidationError(err))
		return
	}

	brand, err := h.brandService.Create(c.Request.Context(), &req)
	if err != nil {
		response.FromAppError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Brand created successfully", brand)
}

// GetByID handles retrieving a brand by ID
func (h *BrandHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	brand, err := h.brandService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.FromAppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Brand retrieved successfully", brand)
}

// GetBySlug handles retrieving a brand by Slug
func (h *BrandHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	brand, err := h.brandService.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		response.FromAppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Brand retrieved successfully", brand)
}

// List handles listing brands with search, is_active/inactive filter, and pagination
func (h *BrandHandler) List(c *gin.Context) {
	p := pagination.GetPagination(c)
	search := c.Query("search")

	var isActive *bool
	if val := c.Query("is_active"); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			isActive = &b
		}
	} else if val := c.Query("inactive"); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			inv := !b
			isActive = &inv
		}
	}

	brands, meta, err := h.brandService.List(c.Request.Context(), search, isActive, p)
	if err != nil {
		response.FromAppError(c, err)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "Brands retrieved successfully", brands, meta)
}

// Update handles updating an existing brand
func (h *BrandHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req application.UpdateBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", validator.FormatValidationError(err))
		return
	}

	brand, err := h.brandService.Update(c.Request.Context(), id, &req)
	if err != nil {
		response.FromAppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Brand updated successfully", brand)
}

// Delete handles deleting a brand
func (h *BrandHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.brandService.Delete(c.Request.Context(), id); err != nil {
		response.FromAppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Brand deleted successfully", nil)
}

// UploadLogo handles uploading a brand logo to MinIO
func (h *BrandHandler) UploadLogo(c *gin.Context) {
	file, err := c.FormFile("logo")
	if err != nil {
		file, err = c.FormFile("file")
		if err != nil {
			response.Error(c, http.StatusBadRequest, "No file uploaded", "File with field name 'logo' or 'file' is required")
			return
		}
	}

	url, err := h.brandService.UploadLogo(c.Request.Context(), file)
	if err != nil {
		response.FromAppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Logo uploaded successfully", gin.H{
		"logo_url": url,
	})
}
