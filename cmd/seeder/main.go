package main

import (
	"context"
	"log"
	"time"

	"github.com/ramdhanrizkij/next-store-api/internal/config"
	"github.com/ramdhanrizkij/next-store-api/internal/database"
	"github.com/ramdhanrizkij/next-store-api/internal/seeder"
)

func main() {
	// 1. Load application and database configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Connecting to database '%s' at %s:%s...", cfg.DB.Name, cfg.DB.Host, cfg.DB.Port)

	// 2. Establish database connection
	db, err := database.NewPostgresDB(&cfg.DB)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Printf("Error closing database connection: %v", err)
			}
		}
	}()

	// 3. Initialize Seeder Registry
	registry := seeder.NewRegistry(db)

	// 4. Execute seeders with timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := registry.RunAll(ctx); err != nil {
		log.Fatalf("Seeding process failed: %v", err)
	}

	log.Println("Database seeding completed successfully.")
}
