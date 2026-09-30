package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ramdhanrizkij/next-store-api/internal/config"
	"github.com/ramdhanrizkij/next-store-api/internal/database"
	"github.com/ramdhanrizkij/next-store-api/internal/server"
)

func main() {
	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting %s in %s mode...", cfg.App.Name, cfg.App.Env)

	// 2. Initialize PostgreSQL Database
	db, err := database.NewPostgresDB(&cfg.DB)
	if err != nil {
		log.Printf("[WARN] Database connection failed: %v", err)
		log.Println("[WARN] Running with limited functionality until database is available")
	} else {
		defer func() {
			if sqlDB, err := db.DB(); err == nil {
				if err := sqlDB.Close(); err != nil {
					log.Printf("Error closing database connection: %v", err)
				}
			}
		}()
		log.Println("Database connection established successfully")
	}

	// 3. Initialize Router & Server
	router := server.NewRouter(cfg, db)
	srv := server.NewServer(cfg, router)

	// 4. Start HTTP Server in background
	go func() {
		log.Printf("HTTP Server is listening on port %s", cfg.App.Port)
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed to listen: %v", err)
		}
	}()

	// 5. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
