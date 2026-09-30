package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *AuthHandler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.GET("/verify-email", handler.VerifyEmail)
		auth.POST("/verify-email", handler.VerifyEmail)
	}
}
