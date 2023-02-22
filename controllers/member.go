package controllers

import (
	"http-server/initializers"
	"http-server/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RetreiveMembers(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")
	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	members := []models.Member{}
	teamID, haveTeam := ctx.GetQuery("team_id")
	keyword, haveKeyword := ctx.GetQuery("keyword")

	var query = initializers.DB.Limit(intLimit).Offset(offset)

	if haveTeam {
		query.Where("team_id = ?", teamID)
	}
	if haveKeyword {
		query.Where("name LIKE ?", "%"+keyword+"%")
	}
	query.Preload("Team").Preload("Roles").Find(&members)

	if query.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": query.Error,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   members,
	})
}

func CreateMember(ctx *gin.Context) {
	var payload *models.CreateMemberRequest

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	roles := []*models.Role{}
	projects := []*models.Project{}

	initializers.DB.Find(&roles, payload.Roles)
	initializers.DB.Find(&projects, payload.Projects)

	member := models.Member{
		Name:      payload.Name,
		Email:     payload.Email,
		Roles:     roles,
		Projects:  projects,
		TeamID:    int(payload.TeamID),
		WorkModel: models.WorkModel(payload.WorkModel),
		Salary:    float32(payload.Salary),
		OtherCost: float32(payload.OtherCost),
		StartDate: payload.StartDate,
	}

	result := initializers.DB.Create(&member)

	if result.Error != nil {
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

func UpdateMember(ctx *gin.Context) {
	id := ctx.Param("id")
	var payload *models.CreateMemberRequest

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	member := models.Member{
		Name:      payload.Name,
		Email:     payload.Email,
		TeamID:    int(payload.TeamID),
		WorkModel: models.WorkModel(payload.WorkModel),
		Salary:    float32(payload.Salary),
		OtherCost: float32(payload.OtherCost),
	}

	result := initializers.DB.Model(&member).Where("id = ?", id).Updates(&member)

	memberId, _ := strconv.Atoi(id)

	if roleErr := AssignRolesToMember(memberId, payload.Roles); roleErr != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": roleErr.Error(),
		})
	}

	if err := AssignProjectsToMember(memberId, payload.Projects); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
	}

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": "Cannot update project",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   member,
	})
}

func GetMember(ctx *gin.Context) {
	id := ctx.Param("id")

	member := models.Member{}
	result := initializers.DB.Where("id = ?", id).Preload("Roles").Preload("Projects").Preload("Team").First(&member, id)

	if result.Error != nil {
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

func DeleteMember(ctx *gin.Context) {
	id := ctx.Param("id")
	result := initializers.DB.Delete(&models.Member{}, id)

	if result.Error != nil {
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

func AssignProjectsToMember(memberId int, projectIds []int) error {
	member := models.Member{}
	projects := []models.Project{}

	initializers.DB.First(&member, memberId)
	initializers.DB.Where("id IN ?", projectIds).Find(&projects)
	err := initializers.DB.Model(&member).Association("Projects").Replace(&projects)

	return err
}

func AssignRolesToMember(memberId int, roleIds []int) error {
	member := models.Member{}
	roles := []models.Role{}

	initializers.DB.First(&member, memberId)
	initializers.DB.Where("id IN ?", roleIds).Find(&roles)
	err := initializers.DB.Model(&member).Association("Roles").Replace(&roles)

	return err
}
