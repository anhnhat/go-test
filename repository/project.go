package repository

import (
	"fmt"
	"strconv"
	"time"

	"http-server/initializers"
	"http-server/models"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ProjectRepo struct {
	DB *gorm.DB
}

var Status = map[string]int{
	"Planning":   0,
	"Inprogress": 1,
	"Completed":  2,
	"Overdue":    3,
	"Cancel":     4,
}

func NewProjectRepo(db *gorm.DB) models.IProjectRepo {
	return &ProjectRepo{
		DB: db,
	}
}

func (pr *ProjectRepo) GetAll(ctx *gin.Context) ([]models.Project, error) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")
	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var projects = []models.Project{}
	status, haveStatus := ctx.GetQuery("status")
	keyword, haveKeyword := ctx.GetQuery("keyword")

	var query = pr.DB.Limit(intLimit).Offset(offset)

	if haveStatus {
		query.Where("status_id = ?", Status[status])
	}
	if haveKeyword {
		query.Where("name LIKE ?", "%"+keyword+"%")
	}
	query.Preload("Members").Find(&projects)

	return projects, query.Error
}

func (pr *ProjectRepo) Create(payload *models.CreateProjectRequest) (models.Project, error) {
	members := []*models.Member{}
	membersResult := pr.DB.Where("id IN ?", payload.Members).Find(&members)

	if membersResult.Error != nil {
		fmt.Println(membersResult.Error)
	}

	project := models.Project{
		Name:           payload.Name,
		Description:    payload.Description,
		StatusID:       payload.StatusID,
		Priority:       payload.Priority,
		Health:         payload.Health,
		ClientName:     payload.ClientName,
		Members:        members,
		Budget:         payload.Budget,
		ActualReceived: payload.ActualReceived,
		StartDate:      payload.StartDate,
		EndDate:        payload.EndDate,
	}

	result := initializers.DB.Create(&project)
	return project, result.Error
}

func (pr *ProjectRepo) GetByIdOrName(id interface{}) (models.Project, error) {
	strId, _ := id.(string)
	var project = models.Project{}
	id, err := strconv.Atoi(strId)

	var queryErr error
	if err != nil {
		queryErr = pr.DB.Where("name LIKE ?", "%"+strId+"%").Preload("Members").First(&project).Error
	} else {
		queryErr = pr.DB.Where("id = ?", id).Preload("Members").First(&project).Error
	}

	return project, queryErr
}

func (pr *ProjectRepo) Update(id int, payload *models.CreateProjectRequest, ctx *gin.Context) (models.Project, error) {
	project := models.Project{
		Name:           payload.Name,
		Description:    payload.Description,
		StatusID:       payload.StatusID,
		Priority:       payload.Priority,
		Health:         payload.Health,
		StartDate:      payload.StartDate,
		EndDate:        payload.EndDate,
		ClientName:     payload.ClientName,
		Budget:         payload.Budget,
		ActualReceived: payload.ActualReceived,
	}

	err := pr.DB.Model(&project).Where("id = ?", id).Updates(&project).Error
	return project, err
}

func (pr *ProjectRepo) Delete(id int) error {
	result := initializers.DB.Delete(&models.Project{}, id)
	// result := initializers.DB.Model(&models.Project{}, id)
	return result.Error
}

func (pr *ProjectRepo) AssignMembersToProject(projectId int, memberIds []uint) error {
	project := models.Project{}
	members := []models.Member{}

	pr.DB.First(&project, projectId)
	pr.DB.Where("id IN ?", memberIds).Find(&members)
	err := pr.DB.Model(&project).Association("Members").Replace(&members)

	timeSheet := NewTimeSheetRepo(pr.DB)
	for _, memberId := range memberIds {
		timeSheet.Create(models.CreateTimeSheetRequest{
			MemberID:  int(memberId),
			ProjectID: projectId,
			DayLog: models.DayLog{
				Date: datatypes.Date(time.Now()),
			},
		})
	}

	return err
}
