package controllers

import (
	"net/http"
	"strconv"

	"http-server/models"

	"github.com/gin-gonic/gin"
)

type MemberController struct {
	MemberUC models.IMemberRepo
}

func (mc *MemberController) RetreiveMembers(ctx *gin.Context) {
	members, err := mc.MemberUC.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   members,
	})
}

func (mc *MemberController) CreateMember(ctx *gin.Context) {
	var payload *models.CreateMemberRequest

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	member, err := mc.MemberUC.CreateMember(payload, ctx)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Cannot create member",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   member,
	})
}

func (mc *MemberController) UpdateMember(ctx *gin.Context) {
	id := ctx.Param("id")
	var payload *models.CreateMemberRequest

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	memberId, _ := strconv.Atoi(id)
	member, err := mc.MemberUC.UpdateMember(memberId, payload, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	mc.MemberUC.AssignRolesToMember(memberId, payload.Roles)
	mc.MemberUC.AssignProjectsToMember(memberId, payload.Projects)

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   member,
	})
}

func (mc *MemberController) GetMember(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.Atoi(id)

	member, err := mc.MemberUC.GetMemberById(idInt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Cannot find member",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   member,
	})
}

func (mc *MemberController) DeleteMember(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.Atoi(id)

	if err := mc.MemberUC.DeleteMember(idInt); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Cannot delete member",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Deleted",
	})
}
