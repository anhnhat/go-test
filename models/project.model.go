package models

import (
	"time"

	"gorm.io/gorm"
)

// TODO
type Priority struct {
}

type Health struct {
	Name   string
	Reason string
}

type Project struct {
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"type:varchar(50);not null"`
	Description   string `gorm:"type:varchar(255)"`
	StatusID      uint
	ProjectStatus ProjectStatus `gorm:"foreignKey:StatusID"`
	Priority      string
	Health        Health  `gorm:"embedded"`
	Users         []*User `gorm:"many2many:user_projects;"`
	StartDate     time.Time
	EndDate       time.Time
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	// HealthReason string  `gorm:"type:varchar(255)"`
}

type ProjectResponse struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	StatusID    uint   `json:"statusId,omitempty"`
	Priority    string `json:"priority,omitempty"`
}
