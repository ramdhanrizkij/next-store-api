package domain

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          string         `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	Name        string         `json:"name" gorm:"type:varchar(100);not null;unique;column:name"`
	Description *string        `json:"description,omitempty" gorm:"type:text;column:description"`
	CreatedAt   time.Time      `json:"created_at" gorm:"type:timestamptz;default:now();column:created_at"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"type:timestamptz;default:now();column:updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;column:deleted_at"`

	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;foreignKey:ID;joinForeignKey:RoleID;References:ID;joinReferences:PermissionID"`
}

func (Role) TableName() string {
	return "roles"
}

// RoleModel alias for backward compatibility
type RoleModel = Role

type Permission struct {
	ID        string         `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	Name      string         `json:"name" gorm:"type:varchar(100);not null;column:name"`
	Module    *string        `json:"module,omitempty" gorm:"type:varchar(100);column:module"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamptz;default:now();column:created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamptz;default:now();column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;column:deleted_at"`
}

func (Permission) TableName() string {
	return "permissions"
}

type UserRole struct {
	ID        string         `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	UserID    string         `json:"user_id" gorm:"type:uuid;not null;column:user_id"`
	RoleID    string         `json:"role_id" gorm:"type:uuid;not null;column:role_id"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamptz;default:now();column:created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamptz;default:now();column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;column:deleted_at"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

type RolePermission struct {
	ID           string         `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	RoleID       string         `json:"role_id" gorm:"type:uuid;not null;column:role_id"`
	PermissionID string         `json:"permission_id" gorm:"type:uuid;not null;column:permission_id"`
	CreatedAt    time.Time      `json:"created_at" gorm:"type:timestamptz;default:now();column:created_at"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"type:timestamptz;default:now();column:updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index;column:deleted_at"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
