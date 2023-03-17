package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name    string    `gorm:"uniqueIndex,not null"`
	Members []*Member `gorm:"many2many:member_roles;"`
}
