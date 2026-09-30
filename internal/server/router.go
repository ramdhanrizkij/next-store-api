package server

import (
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/ramdhanrizkij/next-store-api/internal/config"
	"github.com/ramdhanrizkij/next-store-api/internal/middleware"
	authApp "github.com/ramdhanrizkij/next-store-api/internal/modules/auth/application"
	authHttp "github.com/ramdhanrizkij/next-store-api/internal/modules/auth/delivery/http"
	authInfra "github.com/ramdhanrizkij/next-store-api/internal/modules/auth/infrastructure"
	healthHttp "github.com/ramdhanrizkij/next-store-api/internal/modules/health/delivery/http"
	userApp "github.com/ramdhanrizkij/next-store-api/internal/modules/user/application"
	userHttp "github.com/ramdhanrizkij/next-store-api/internal/modules/user/delivery/http"
	userInfra "github.com/ramdhanrizkij/next-store-api/internal/modules/user/infrastructure"
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

	// API v1 group
	v1 := r.Group("/api/v1")
	{
		// Dependencies initialization
		userRepo := userInfra.NewUserPostgresRepository(db)
		authRepo := authInfra.NewAuthPostgresRepository(db)

		// Queue / Worker Task Distributor
		redisOpt := asynq.RedisClientOpt{
			Addr:     cfg.Redis.Addr(),
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		}
		taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

		userService := userApp.NewUserService(userRepo)
		authService := authApp.NewAuthService(
			userRepo,
			authRepo,
			taskDistributor,
			cfg.App.URL,
			cfg.JWT.Secret,
			cfg.JWT.ExpirationHours,
		)

		userHandler := userHttp.NewUserHandler(userService)
		authHandler := authHttp.NewAuthHandler(authService)

		authMiddleware := middleware.Auth(cfg.JWT.Secret)

		// Module Routes
		authHttp.RegisterRoutes(v1, authHandler)
		userHttp.RegisterRoutes(v1, userHandler, authMiddleware)
	}

	return r
}
