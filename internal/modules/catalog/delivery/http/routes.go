package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, brandHandler *BrandHandler, authMiddleware gin.HandlerFunc) {
	brands := rg.Group("/brands")
	{
		// Public routes
		brands.GET("", brandHandler.List)
		brands.GET("/:id", brandHandler.GetByID)
		brands.GET("/slug/:slug", brandHandler.GetBySlug)

		// Protected routes
		protected := brands.Group("")
		if authMiddleware != nil {
			protected.Use(authMiddleware)
		}
		{
			protected.POST("", brandHandler.Create)
			protected.PUT("/:id", brandHandler.Update)
			protected.DELETE("/:id", brandHandler.Delete)
			protected.POST("/upload-logo", brandHandler.UploadLogo)
		}
	}
}
