package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ramdhanrizkij/next-store-api/internal/shared/response"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC RECOVER] %v\n", r)
				response.Error(c, http.StatusInternalServerError, "Internal Server Error", "An unexpected error occurred")
				c.Abort()
			}
		}()
		c.Next()
	}
}
