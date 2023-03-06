package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type TimeSheetSegment struct {
	gorm.Model
	TimeSheetID uint `gorm:"not_null"`
	TimeSheet   TimeSheet
	Hours       float32
	Date        datatypes.Date `gorm:"default:NULL"`
}

type ITimeSheetSegment interface {
	Create()
}
