package models

import (
	"time"
)

type Role struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Users     []*User   `gorm:"many2many:user_roles;"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time
}
