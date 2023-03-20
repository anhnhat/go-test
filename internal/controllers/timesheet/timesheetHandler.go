package timesheet

import (
	"net/http"
	"strconv"

	"http-server/internal/appctx"
	"http-server/internal/models"
	"http-server/internal/services/timesheet"

	"github.com/gin-gonic/gin"
)

type timesheetHandler struct {
	appCtx           *appctx.AppCtx
	timesheetService timesheet.ITimeSheet
}

func NewTimesheetHandler(appCtx *appctx.AppCtx) timesheetHandler {
	return timesheetHandler{
		appCtx:           appCtx,
		timesheetService: appCtx.GetTimesheetService(),
	}
}

func (th *timesheetHandler) GetAll(ctx *gin.Context) {
	from, _ := ctx.GetQuery("from")
	timeSheets := th.timesheetService.GetAll(from)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   timeSheets,
	})
}

func (th *timesheetHandler) GetByMemberId(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.Atoi(id)
	th.timesheetService.GetByMemberId(idInt, ctx)
}

func (th *timesheetHandler) Create(ctx *gin.Context) {
	var payload *models.CreateTimeSheetRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timesheet := models.TimeSheet{
		MemberID:  payload.MemberID,
		ProjectID: payload.ProjectID,
		DayLog:    payload.DayLog,
	}

	if err := th.timesheetService.Create(timesheet); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": payload})
}

func (th *timesheetHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.Atoi(id)
	if err := th.timesheetService.Delete(idInt); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Delete timesheet successful",
	})
}

func (th *timesheetHandler) SaveSegment(ctx *gin.Context) {
	timesheetId := ctx.Param("id")
	timesheetIdInt, _ := strconv.Atoi(timesheetId)
	var tsSegmentsPayload []models.CreateTimesheetSegment
	if err := ctx.ShouldBindJSON(&tsSegmentsPayload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var segmentModels []models.TimeSheetSegment
	for _, segment := range tsSegmentsPayload {
		segmentModels = append(segmentModels, models.TimeSheetSegment{
			TimeSheetID: uint(timesheetIdInt),
			Hours:       segment.Hours,
			Date:        segment.Date,
		})
	}
	th.timesheetService.SaveSegment(timesheetIdInt, segmentModels)
}

func (th *timesheetHandler) SaveMultipleTsSegment(ctx *gin.Context) {
	var tsSegmentsPayload []models.MultipleTsSegmentRequest

	if err := ctx.ShouldBindJSON(&tsSegmentsPayload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// TODO: Handle error on save
	errFlag := false // not work
	for _, segment := range tsSegmentsPayload {
		var segmentModels []models.TimeSheetSegment
		for _, segmentNode := range segment.Segments {
			segmentModels = append(segmentModels, models.TimeSheetSegment{
				TimeSheetID: uint(segment.TimesheetId),
				Hours:       segmentNode.Hours,
				Date:        segmentNode.Date,
			})
		}
		err := th.timesheetService.SaveSegment(segment.TimesheetId, segmentModels)
		if err != nil {
			errFlag = true
			ctx.JSON(http.StatusBadGateway, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		}
	}
	if errFlag == false {
		ctx.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"message": "Timesheet created successful",
		})
	}
}

func (th *timesheetHandler) GetById(ctx *gin.Context) {
	timesheetId, _ := strconv.Atoi(ctx.Param("id"))
	timesheet, _ := th.timesheetService.GetById(timesheetId)
	ctx.JSON(http.StatusOK, gin.H{"data": timesheet})
}
