package models

import (
	"time"
)

type ProjectStatus struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(20);not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time
}
