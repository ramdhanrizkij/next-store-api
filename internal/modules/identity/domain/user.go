package domain

import (
	"time"

	"gorm.io/gorm"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

type RoleName = string

const (
	RoleUser  RoleName = "user"
	RoleAdmin RoleName = "admin"
)

type User struct {
	ID              string         `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	Name            string         `json:"name" gorm:"type:varchar(150);not null;column:name"`
	Email           string         `json:"email" gorm:"type:varchar(255);not null;uniqueIndex;column:email"`
	Password        string         `json:"-" gorm:"type:varchar(255);not null;column:password"`
	Phone           *string        `json:"phone,omitempty" gorm:"type:varchar(50);column:phone"`
	Status          UserStatus     `json:"status" gorm:"type:user_status;not null;default:'inactive';column:status"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at,omitempty" gorm:"type:timestamptz;column:email_verified_at"`
	LastLoginAt     *time.Time     `json:"last_login_at,omitempty" gorm:"type:timestamptz;column:last_login_at"`
	CreatedAt       time.Time      `json:"created_at" gorm:"type:timestamptz;default:now();column:created_at"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"type:timestamptz;default:now();column:updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index;column:deleted_at"`

	// Relationships
	Roles []Role `json:"roles,omitempty" gorm:"many2many:user_roles;foreignKey:ID;joinForeignKey:UserID;References:ID;joinReferences:RoleID"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) IsVerified() bool {
	return u.EmailVerifiedAt != nil
}

func (u *User) PrimaryRole() string {
	if len(u.Roles) > 0 {
		return u.Roles[0].Name
	}
	return RoleUser
}
