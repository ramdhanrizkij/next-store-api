package http

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ramdhanrizkij/next-store-api/internal/shared/response"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
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
		if err := h.db.PingContext(ctx); err != nil {
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
