package models

import (
	"time"

	"gorm.io/gorm"
)

type WorkModel int

const (
	Onsite WorkModel = iota
	Remote
	Hybrid
)

type Member struct {
	gorm.Model
	Name      string     `gorm:"type:varchar(255);not null"`
	Email     string     `gorm:"uniqueIndex;not null"`
	Roles     []*Role    `gorm:"many2many:member_roles;"`
	Projects  []*Project `gorm:"many2many:member_projects;"`
	TeamID    int        `gorm:"default:null"`
	Team      Team       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	WorkModel WorkModel  `gorm:"default:0"`
	Salary    float32    `gorm:"type:decimal(20,3);default:null;"`
	OtherCost float32    `gorm:"type:decimal(20,3);default:null;"`
	StartDate time.Time  `gorm:"default:NULL"`
}

type CreateMemberRequest struct {
	Name      string    `json:"Name" binding:"required"`
	Email     string    `json:"Email,omitempty" binding:"required"`
	TeamID    uint      `json:"TeamID,omitempty"`
	Roles     []int     `json:"Roles"`
	WorkModel int       `json:"WorkModel"`
	Salary    float32   `json:"Salary,omitempty"`
	OtherCost float32   `json:"OtherCost,omitempty"`
	Projects  []int     `json:"ProjectIds,omitempty"`
	StartDate time.Time `json:"StartDate,omitempty"`
}
