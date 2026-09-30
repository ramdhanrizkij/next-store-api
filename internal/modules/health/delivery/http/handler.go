package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ramdhanrizkij/next-store-api/internal/shared/response"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Ping(c *gin.Context) {
	response.Success(c, http.StatusOK, "Pong", gin.H{
		"status": "healthy",
	})
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	dbStatus := "connected"
	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err != nil || sqlDB.PingContext(ctx) != nil {
			dbStatus = "disconnected"
		}
	} else {
		dbStatus = "not configured"
	}

	response.Success(c, http.StatusOK, "Service is healthy", gin.H{
		"app":       "next-store-api",
		"status":    "running",
		"database":  dbStatus,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
