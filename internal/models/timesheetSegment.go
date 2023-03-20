package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type TimeSheetSegment struct {
	gorm.Model
	TimeSheetID uint `gorm:"not_null;uniqueIndex:timesheet_date"`
	// TimeSheet   TimeSheet
	Hours float32
	Date  datatypes.Date `gorm:"default:NULL;uniqueIndex:timesheet_date"`
}

type CreateTimesheetSegment struct {
	Hours float32        `json:"Hours,omitempty"`
	Date  datatypes.Date `json:"Date,omitempty"`
}

type SegmentRequest struct {
	Hours float32        `json:"Hours,omitempty"`
	Date  datatypes.Date `json:"Date,omitempty"`
}

type MultipleTsSegmentRequest struct {
	TimesheetId int `json:"TimesheetId"`
	Segments    []SegmentRequest
}

type ITimeSheetSegment interface {
	Create()
}
