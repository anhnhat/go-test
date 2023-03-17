package timesheet

import (
	"fmt"
	"net/http"
	"time"

	"http-server/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ITimeSheet interface {
	GetAll() []models.TimeSheet
	Create(models.TimeSheet) error
	Update()
	Delete(id int) error
	SaveSegment(int, []models.TimeSheetSegment)
	GetById(id int) (models.TimeSheet, error)
	GetByMemberId(memberId int, ctx *gin.Context)
	GetByProjectId()
}

type TimeSheet struct {
	DB *gorm.DB
}

type TimeSheetResponse struct {
	ProjectName string    `json:"project_id"`
	MemberName  string    `json:"member_id"`
	Date        time.Time `json:"date"`
	TrackTime   int       `json:"track_time"`
}

type FormatTimeSheet struct {
	ProjectName string
	MemberName  string
	Logs        []models.DayLog
}

func NewTimeSheetService(db *gorm.DB) ITimeSheet {
	return &TimeSheet{
		DB: db,
	}
}

func (ts *TimeSheet) GetAll() []models.TimeSheet {
	timeSheets := []models.TimeSheet{}
	result := ts.DB.Limit(10).Offset(0).Preload("Member").Preload("Project").Preload("TimeSheetSegment").Find(&timeSheets)

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

func (ts *TimeSheet) Create(timeSheet models.TimeSheet) error {
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

func (ts *TimeSheet) SaveSegment(timesheetId int, timesheetSegments []models.TimeSheetSegment) {
	_, err := ts.GetById(timesheetId)
	if err != nil {
		fmt.Println("Cannot get timesheet by id")
		return
	}

	tsErr := ts.DB.Create(&timesheetSegments)
	if tsErr.Error != nil {
		fmt.Printf("Error create segment: %s", tsErr.Error.Error())
	}
}

func (ts *TimeSheet) GetById(id int) (models.TimeSheet, error) {
	timesheet := models.TimeSheet{}
	result := ts.DB.Where("id = ?", uint(id)).Preload("Member").Preload("Project").Preload("TimeSheetSegment").Find(&timesheet)
	return timesheet, result.Error
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
