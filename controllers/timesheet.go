package controllers

import (
	"net/http"
	"strconv"

	"http-server/models"

	"github.com/gin-gonic/gin"
)

type TimeSheetController struct {
	TimeSheetUC models.ITimeSheet
}

func (tc *TimeSheetController) GetAll(ctx *gin.Context) {
	timeSheets := tc.TimeSheetUC.GetAll()
	ctx.JSON(http.StatusOK, gin.H{"data": timeSheets})
}

func (tc *TimeSheetController) GetByMemberId(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.Atoi(id)
	tc.TimeSheetUC.GetByMemberId(idInt, ctx)
}

func (tc *TimeSheetController) Create(ctx *gin.Context) {
	var payload models.CreateTimeSheetRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if err := tc.TimeSheetUC.Create(payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": payload})
}

func (tc *TimeSheetController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.Atoi(id)
	tc.TimeSheetUC.Delete(idInt)
}
