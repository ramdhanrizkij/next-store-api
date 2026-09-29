package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *UserHandler, authMiddleware gin.HandlerFunc) {
	users := rg.Group("/users")
	users.Use(authMiddleware)
	{
		users.GET("/me", handler.GetMe)
		users.PUT("/me", handler.Update)
		users.GET("", handler.List)
		users.GET("/:id", handler.GetByID)
		users.DELETE("/:id", handler.Delete)
	}
}
