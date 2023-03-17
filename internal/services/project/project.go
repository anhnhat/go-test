package project

import (
	"fmt"
	"strconv"

	"http-server/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IProjectService interface {
	GetAll(ctx *gin.Context) ([]models.Project, error)
	GetByIdOrName(idOrName interface{}) (models.Project, error)
	Create(payload *models.CreateProjectRequest) (models.Project, error)
	Update(id int, payload *models.CreateProjectRequest, ctx *gin.Context) (models.Project, error)
	Delete(id int) error
	AssignMembersToProject(projectId int, memberIds []uint) error
}

type ProjectService struct {
	DB *gorm.DB
}

var Status = map[string]int{
	"Planning":   0,
	"Inprogress": 1,
	"Completed":  2,
	"Overdue":    3,
	"Cancel":     4,
}

func NewProjectService(db *gorm.DB) IProjectService {
	return &ProjectService{
		DB: db,
	}
}

func (ps *ProjectService) GetAll(ctx *gin.Context) ([]models.Project, error) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")
	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var projects = []models.Project{}
	status, haveStatus := ctx.GetQuery("status")
	keyword, haveKeyword := ctx.GetQuery("keyword")

	var query = ps.DB.Limit(intLimit).Offset(offset)

	if haveStatus {
		query.Where("status_id = ?", Status[status])
	}
	if haveKeyword {
		query.Where("name LIKE ?", "%"+keyword+"%")
	}
	query.Preload("Members").Find(&projects)

	return projects, query.Error
}

func (ps *ProjectService) Create(payload *models.CreateProjectRequest) (models.Project, error) {
	members := []*models.Member{}
	membersResult := ps.DB.Where("id IN ?", payload.Members).Find(&members)

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

	result := ps.DB.Create(&project)
	return project, result.Error
}

func (ps *ProjectService) GetByIdOrName(id interface{}) (models.Project, error) {
	strId, _ := id.(string)
	var project = models.Project{}
	id, err := strconv.Atoi(strId)

	var queryErr error
	if err != nil {
		queryErr = ps.DB.Where("name LIKE ?", "%"+strId+"%").Preload("Members").First(&project).Error
	} else {
		queryErr = ps.DB.Where("id = ?", id).Preload("Members").First(&project).Error
	}

	return project, queryErr
}

func (ps *ProjectService) Update(id int, payload *models.CreateProjectRequest, ctx *gin.Context) (models.Project, error) {
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

	err := ps.DB.Model(&project).Where("id = ?", id).Updates(&project).Error
	return project, err
}

func (ps *ProjectService) Delete(id int) error {
	result := ps.DB.Delete(&models.Project{}, id)
	return result.Error
}

func (ps *ProjectService) AssignMembersToProject(projectId int, memberIds []uint) error {
	project := models.Project{}
	members := []models.Member{}

	ps.DB.First(&project, projectId)
	ps.DB.Where("id IN ?", memberIds).Find(&members)
	err := ps.DB.Model(&project).Association("Members").Replace(&members)

	// timeSheet := NewTimeSheetRepo(ps.DB)
	// for _, memberId := range memberIds {
	// 	timeSheet.Create(models.CreateTimeSheetRequest{
	// 		MemberID:  int(memberId),
	// 		ProjectID: projectId,
	// 		DayLog: models.DayLog{
	// 			Date: datatypes.Date(time.Now()),
	// 		},
	// 	})
	// }

	return err
}
