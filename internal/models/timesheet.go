package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DayLog struct {
	Date      datatypes.Date `json:"Date" gorm:"uniqueIndex:mb_prj_date"`
	TrackTime float32        `json:"TrackTime,omitempty"`
}

type TimeSheet struct {
	gorm.Model
	MemberID         int     `gorm:"default:null;uniqueIndex:mb_prj_date"`
	Member           Member  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ProjectID        int     `gorm:"default:null;uniqueIndex:mb_prj_date"`
	Project          Project `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	DayLog           DayLog  `gorm:"embedded;uniqueIndex:mb_prj_date"`
	Note             string  `gorm:"varchar(255)"`
	TimeSheetSegment []TimeSheetSegment
	StartDate        time.Time `gorm:"default:NULL"`
	EndDate          time.Time `gorm:"default:NULL"`
}

type CreateTimeSheetRequest struct {
	MemberID  int `json:"MemberId"`
	ProjectID int `json:"ProjectId"`
	DayLog    DayLog
}

type TimeSheetResponse struct {
	MemberID  int `json:"MemberId"`
	ProjectID int `json:"ProjectId"`
	DayLog    DayLog
}

type ITimeSheet interface {
	GetAll() []TimeSheet
	Create(CreateTimeSheetRequest) error
	Update()
	Delete(id int) error
	GetByMemberId(memberId int, ctx *gin.Context)
	GetByProjectId()
}
