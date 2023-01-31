package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name      string     `gorm:"type:varchar(255);not null"`
	Email     string     `gorm:"uniqueIndex;not null"`
	Password  string     `gorm:"not null"`
	Projects  []*Project `gorm:"many2many:user_projects;"`
	Roles     []*Role    `gorm:"many2many:user_roles;"`
	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt time.Time  `gorm:"not null"`
}

type CreateUserRequest struct {
	Name     string `json:"name"  binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
