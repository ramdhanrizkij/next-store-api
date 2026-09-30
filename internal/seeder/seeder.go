package seeder

import (
	"context"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Seeder interface {
	Name() string
	Seed(ctx context.Context, db *gorm.DB) error
}

type Registry struct {
	db      *gorm.DB
	seeders []Seeder
}

func NewRegistry(db *gorm.DB) *Registry {
	r := &Registry{
		db:      db,
		seeders: make([]Seeder, 0),
	}

	// Register default seeders
	r.Register(NewUserSeeder())

	return r
}

func (r *Registry) Register(s Seeder) {
	r.seeders = append(r.seeders, s)
}

func (r *Registry) RunAll(ctx context.Context) error {
	log.Printf("Starting database seeding (%d seeders registered)...", len(r.seeders))

	for _, s := range r.seeders {
		log.Printf("==> Executing: %s", s.Name())
		if err := s.Seed(ctx, r.db); err != nil {
			return fmt.Errorf("seeder '%s' encountered an error: %w", s.Name(), err)
		}
		log.Printf("==> Completed: %s", s.Name())
	}

	log.Println("All database seeders completed successfully.")
	return nil
}
