package domain

import (
	"time"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID         string    `json:"id" gorm:"primaryKey;column:id"`
	Name       string    `json:"name" gorm:"column:name"`
	Email      string    `json:"email" gorm:"column:email;uniqueIndex"`
	Password   string    `json:"-" gorm:"column:password"`
	Role       Role      `json:"role" gorm:"column:role"`
	IsVerified bool      `json:"is_verified" gorm:"column:is_verified"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "users"
}
