package repository

import (
	"fmt"
	"net/http"
	"time"

	"http-server/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TimeSheet struct {
	DB *gorm.DB
}

type TimeSheetResponse struct {
	ProjectName string    `json:"project_id"`
	MemberName  string    `json:"member_id"`
	Date        time.Time `json:"date"`
	TrackTime   int       `json:"track_time"`
	// DayLog    models.DayLog
}

type FormatTimeSheet struct {
	ProjectName string
	MemberName  string
	Logs        []models.DayLog
}

func NewTimeSheetRepo(db *gorm.DB) models.ITimeSheet {
	return &TimeSheet{
		DB: db,
	}
}

func (ts *TimeSheet) GetAll() []models.TimeSheet {
	timeSheets := []models.TimeSheet{}
	result := ts.DB.Limit(10).Offset(0).Preload("Member").Preload("Project").Find(&timeSheets)

	// formatTimesheets := []interface{}{}
	// logs := []models.DayLog{}
	// for _, timeSheet := range timeSheets {
	// 	logs = append(logs, models.DayLog{
	// 		Date:      timeSheet.DayLog.Date,
	// 		TrackTime: timeSheet.DayLog.TrackTime,
	// 	})
	// 	formatTimesheets = append(formatTimesheets, FormatTimeSheet{
	// 		ProjectName: timeSheet.Project.Name,
	// 		MemberName:  timeSheet.Member.Name,
	// 	})
	// }

	if result.Error != nil {
		fmt.Println("Error on get all timesheet")
	}
	return timeSheets
}

func (ts *TimeSheet) Create(payload models.CreateTimeSheetRequest) error {
	timeSheet := models.TimeSheet{
		MemberID:  payload.MemberID,
		ProjectID: payload.ProjectID,
		DayLog:    payload.DayLog,
	}

	result := ts.DB.Create(&timeSheet)
	if result.Error != nil {
		fmt.Println("Error on create timesheet")
		return result.Error
	}

	return nil
}

func (ts *TimeSheet) Update() {}
func (ts *TimeSheet) Delete(id int) error {
	err := ts.DB.Delete(&models.TimeSheet{}, id).Error
	return err
}

func (ts *TimeSheet) GetByMemberId(userId int, ctx *gin.Context) {
	from := ctx.Query("from")
	// timeSheets := []models.TimeSheet{}
	timeSheets := []map[string]interface{}{}

	// query := ts.DB.Model(&timeSheets).Select("members.name, projects.name, time_sheets.date")
	// if from != "" {
	// 	query = query.Where("date >= ?", from)
	// }
	// result := query.Joins("JOIN members ON members.id = time_sheets.member_id").Joins("JOIN projects ON projects.id = time_sheet.project_id").Where("member_id = ? AND date >= ?", userId, from).Find(&timeSheets)

	result := ts.DB.Raw("SELECT members.name as MemberName, projects.name as ProjectName, time_sheets.date, time_sheets.track_time FROM time_sheets INNER JOIN members ON time_sheets.member_id = members.id INNER JOIN projects ON time_sheets.project_id = projects.id WHERE member_id = ? AND time_sheets.date >= ?", userId, from).Find(&timeSheets)

	// formatedRes := []TimeSheetResponse{}
	// for _, record := range timeSheets {

	// 	projectName := record["projectname"].(string)
	// 	memberName := record["membername"].(string)
	// 	date := record["date"]
	// 	trackTime := record["track_time"].(int)

	// 	formatedRes = append(formatedRes, TimeSheetResponse{
	// 		ProjectName: projectName,
	// 		MemberName:  memberName,
	// 		Date:        date,
	// 		TrackTime:   trackTime,
	// 	})
	// }

	if result.Error != nil {
		fmt.Println("Cannot get timesheet")
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": timeSheets,
	})
}
func (ts *TimeSheet) GetByProjectId() {}
