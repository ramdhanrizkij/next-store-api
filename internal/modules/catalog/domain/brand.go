package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Brand struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(150);not null" json:"name"`
	Slug        string         `gorm:"type:varchar(180);not null;unique" json:"slug"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	LogoUrl     *string        `gorm:"type:text" json:"logo_url,omitempty"`
	IsActive    bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt   time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Brand) TableName() string {
	return "brands"
}
