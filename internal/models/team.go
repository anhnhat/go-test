package models

import "gorm.io/gorm"

type Team struct {
	gorm.Model
	Name string `gorm:"type:varchar(255);not null"`
}
