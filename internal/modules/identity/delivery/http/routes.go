package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, authHandler *AuthHandler, userHandler *UserHandler, authMiddleware gin.HandlerFunc) {
	// Authentication routes
	auth := rg.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/verify-email", authHandler.VerifyEmail)
		auth.POST("/verify-email", authHandler.VerifyEmail)
	}

	// User management routes (protected)
	users := rg.Group("/users")
	users.Use(authMiddleware)
	{
		users.GET("/me", userHandler.GetMe)
		users.PUT("/me", userHandler.Update)
		users.GET("", userHandler.List)
		users.GET("/:id", userHandler.GetByID)
		users.DELETE("/:id", userHandler.Delete)
	}
}
