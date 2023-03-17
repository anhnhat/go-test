package member

import (
	"net/http"
	"strconv"

	"http-server/internal/appctx"
	"http-server/internal/models"
	"http-server/internal/services/member"

	"github.com/gin-gonic/gin"
)

type memberHandler struct {
	appctx        *appctx.AppCtx
	memberService member.IMemberService
}

func NewMemberHandler(appCtx *appctx.AppCtx) memberHandler {
	return memberHandler{
		appctx:        appCtx,
		memberService: appCtx.GetMemberService(),
	}
}

func (mh *memberHandler) RetreiveMembers(ctx *gin.Context) {
	members, err := mh.memberService.GetAll(ctx)
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

func (mh *memberHandler) CreateMember(ctx *gin.Context) {
	var payload *models.CreateMemberRequest

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	member, err := mh.memberService.CreateMember(payload, ctx)

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

func (mh *memberHandler) UpdateMember(ctx *gin.Context) {
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
	member, err := mh.memberService.UpdateMember(memberId, payload, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	mh.memberService.AssignRolesToMember(memberId, payload.Roles)
	mh.memberService.AssignProjectsToMember(memberId, payload.Projects)

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   member,
	})
}

func (mh *memberHandler) GetMember(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.Atoi(id)

	member, err := mh.memberService.GetMemberById(idInt)
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

func (mh *memberHandler) DeleteMember(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.Atoi(id)

	if err := mh.memberService.DeleteMember(idInt); err != nil {
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
