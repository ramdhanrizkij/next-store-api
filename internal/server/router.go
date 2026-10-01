package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/ramdhanrizkij/next-store-api/internal/config"
	"github.com/ramdhanrizkij/next-store-api/internal/middleware"
	catalogApp "github.com/ramdhanrizkij/next-store-api/internal/modules/catalog/application"
	catalogHttp "github.com/ramdhanrizkij/next-store-api/internal/modules/catalog/delivery/http"
	catalogInfra "github.com/ramdhanrizkij/next-store-api/internal/modules/catalog/infrastructure"
	healthHttp "github.com/ramdhanrizkij/next-store-api/internal/modules/health/delivery/http"
	identityApp "github.com/ramdhanrizkij/next-store-api/internal/modules/identity/application"
	identityHttp "github.com/ramdhanrizkij/next-store-api/internal/modules/identity/delivery/http"
	identityInfra "github.com/ramdhanrizkij/next-store-api/internal/modules/identity/infrastructure"
	"github.com/ramdhanrizkij/next-store-api/internal/shared/storage"
	"github.com/ramdhanrizkij/next-store-api/internal/worker"
	"gorm.io/gorm"
)

func NewRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Global Middlewares
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// Health Check Route
	healthHandler := healthHttp.NewHealthHandler(db)
	r.GET("/health", healthHandler.HealthCheck)
	r.GET("/ping", healthHandler.Ping)

	// Object Storage Service (MinIO)
	storageService, err := storage.NewMinioStorageService(cfg.Minio)
	if err != nil {
		log.Printf("[WARN] Failed to initialize MinIO storage: %v", err)
	}

	// API v1 group
	v1 := r.Group("/api/v1")
	{
		// Dependencies initialization (Identity Domain)
		userRepo := identityInfra.NewUserPostgresRepository(db)
		authRepo := identityInfra.NewAuthPostgresRepository(db)

		// Queue / Worker Task Distributor
		redisOpt := asynq.RedisClientOpt{
			Addr:     cfg.Redis.Addr(),
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		}
		taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

		userService := identityApp.NewUserService(userRepo)
		authService := identityApp.NewAuthService(
			userRepo,
			authRepo,
			taskDistributor,
			cfg.App.URL,
			cfg.JWT.Secret,
			cfg.JWT.ExpirationHours,
		)

		userHandler := identityHttp.NewUserHandler(userService)
		authHandler := identityHttp.NewAuthHandler(authService)

		authMiddleware := middleware.Auth(cfg.JWT.Secret)

		// Register Identity domain routes
		identityHttp.RegisterRoutes(v1, authHandler, userHandler, authMiddleware)

		// Dependencies initialization (Catalog Domain)
		brandRepo := catalogInfra.NewBrandPostgresRepository(db)
		brandService := catalogApp.NewBrandService(brandRepo, storageService)
		brandHandler := catalogHttp.NewBrandHandler(brandService)

		// Register Catalog domain routes
		catalogHttp.RegisterRoutes(v1, brandHandler, authMiddleware)
	}

	return r
}
