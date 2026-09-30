package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/ramdhanrizkij/next-store-api/internal/config"
	"github.com/ramdhanrizkij/next-store-api/internal/worker"
)

func main() {
	// 1. Load application configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting %s background worker...", cfg.App.Name)

	// 2. Configure Redis connection for Asynq
	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.Redis.Addr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}

	// 3. Initialize Mailer and Task Processor
	mailer := worker.NewLogMailer()
	processor := worker.NewRedisTaskProcessor(redisOpt, mailer)

	// 4. Start worker server in a background goroutine
	go func() {
		if err := processor.Start(); err != nil {
			log.Fatalf("Worker server failed: %v", err)
		}
	}()

	// 5. Wait for termination signals (Graceful Shutdown)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker gracefully...")
	processor.Shutdown()
	log.Println("Worker exited cleanly.")
}
