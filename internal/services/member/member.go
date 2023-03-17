package member

import (
	"strconv"

	"http-server/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IMemberService interface {
	GetAll(ctx *gin.Context) ([]models.Member, error)
	CreateMember(payload *models.CreateMemberRequest, ctx *gin.Context) (models.Member, error)
	UpdateMember(id int, payload *models.CreateMemberRequest, ctx *gin.Context) (models.Member, error)
	GetMemberById(id int) (models.Member, error)
	DeleteMember(id int) error
	AssignProjectsToMember(memberId int, projectIds []int) error
	AssignRolesToMember(memberId int, roleIds []int) error
}

type MemberService struct {
	DB *gorm.DB
}

func NewMemberService(db *gorm.DB) models.IMemberService {
	return &MemberService{
		DB: db,
	}
}

func (mr *MemberService) GetAll(ctx *gin.Context) ([]models.Member, error) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")
	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	teamID, haveTeam := ctx.GetQuery("team_id")
	keyword, haveKeyword := ctx.GetQuery("keyword")
	members := []models.Member{}

	var query = mr.DB.Limit(intLimit).Offset(offset)
	if haveTeam {
		query.Where("team_id = ?", teamID)
	}
	if haveKeyword {
		query.Where("name LIKE ?", "%"+keyword+"%")
	}
	err := query.Preload("Team").Preload("Roles").Find(&members).Error

	return members, err
}

func (mr *MemberService) CreateMember(payload *models.CreateMemberRequest, ctx *gin.Context) (models.Member, error) {
	roles := []*models.Role{}
	projects := []*models.Project{}

	mr.DB.Find(&roles, payload.Roles)
	mr.DB.Find(&projects, payload.Projects)

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

	err := mr.DB.Create(&member).Error
	return member, err
}

func (mr *MemberService) UpdateMember(id int, payload *models.CreateMemberRequest, ctx *gin.Context) (models.Member, error) {
	member := models.Member{
		Name:      payload.Name,
		Email:     payload.Email,
		TeamID:    int(payload.TeamID),
		WorkModel: models.WorkModel(payload.WorkModel),
		Salary:    float32(payload.Salary),
		OtherCost: float32(payload.OtherCost),
	}

	err := mr.DB.Model(&member).Where("id = ?", id).Updates(&member).Error
	return member, err
}

func (mr *MemberService) GetMemberById(id int) (models.Member, error) {
	member := models.Member{}
	result := mr.DB.Where("id = ?", id).Preload("Roles").Preload("Projects").Preload("Team").First(&member, id)
	return member, result.Error
}

func (mr *MemberService) DeleteMember(id int) error {
	result := mr.DB.Delete(&models.Member{}, id).Error
	return result
}

func (mr *MemberService) AssignProjectsToMember(memberId int, projectIds []int) error {
	member := models.Member{}
	projects := []models.Project{}

	mr.DB.First(&member, memberId)
	mr.DB.Where("id IN ?", projectIds).Find(&projects)
	err := mr.DB.Model(&member).Association("Projects").Replace(&projects)

	// timeSheet := NewTimeSheetRepo(mr.DB)
	// for _, projectId := range projectIds {
	// 	timeSheet.Create(models.CreateTimeSheetRequest{
	// 		MemberID:  memberId,
	// 		ProjectID: projectId,
	// 		DayLog: models.DayLog{
	// 			Date: datatypes.Date(time.Now()),
	// 		},
	// 	})
	// }

	return err
}

func (mr *MemberService) AssignRolesToMember(memberId int, roleIds []int) error {
	member := models.Member{}
	roles := []models.Role{}

	mr.DB.First(&member, memberId)
	mr.DB.Where("id IN ?", roleIds).Find(&roles)
	err := mr.DB.Model(&member).Association("Roles").Replace(&roles)

	return err
}
