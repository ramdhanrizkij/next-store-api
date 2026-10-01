package pagination

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}

type Meta struct {
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

func GetPagination(c *gin.Context) Pagination {
	page, err := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(DefaultPage)))
	if err != nil || page < 1 {
		page = DefaultPage
	}

	limitStr := c.Query("pageSize")
	if limitStr == "" {
		limitStr = c.Query("page_size")
	}
	if limitStr == "" {
		limitStr = c.DefaultQuery("limit", strconv.Itoa(DefaultLimit))
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	return Pagination{
		Page:  page,
		Limit: limit,
	}
}

func BuildMeta(totalItems int64, page, limit int) *Meta {
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	return &Meta{
		CurrentPage: page,
		Limit:       limit,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
	}
}
