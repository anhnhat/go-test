package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key"`
	Name       string    `gorm:"type:varchar(255);not null"`
	Email      string    `gorm:"uniqueIndex;not null"`
	Password   string    `gorm:"not null"`
	AllowLogin bool      `gorm:"default:true"`
	Roles      []*Role   `gorm:"many2many:user_roles;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CreateUserRequest struct {
	Name     string `json:"Name"  binding:"required"`
	Email    string `json:"Email" binding:"required"`
	Password string `json:"Password" binding:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"ID,omitempty"`
	Name      string    `json:"Name,omitempty"`
	Email     string    `json:"Email,omitempty"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
