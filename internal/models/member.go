package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkModel int

const (
	Onsite WorkModel = iota
	Remote
	Hybrid
)

type Member struct {
	gorm.Model
	Name      string     `gorm:"index:idx_members_email;type:varchar(255);not null"`
	Email     string     `gorm:"index:idx_members_email;not null"`
	Roles     []*Role    `gorm:"many2many:member_roles;"`
	Projects  []*Project `gorm:"many2many:member_projects;"`
	TeamID    int        `gorm:"default:null"`
	Team      Team       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	WorkModel WorkModel  `gorm:"default:0"`
	Salary    float32    `gorm:"type:decimal(20,3);default:null;"`
	OtherCost float32    `gorm:"type:decimal(20,3);default:null;"`
	StartDate time.Time  `gorm:"default:NULL"`
}

type CreateMemberRequest struct {
	Name      string    `json:"Name" binding:"required"`
	Email     string    `json:"Email,omitempty" binding:"required"`
	TeamID    uint      `json:"TeamID,omitempty"`
	Roles     []int     `json:"Roles"`
	WorkModel int       `json:"WorkModel"`
	Salary    float32   `json:"Salary,omitempty"`
	OtherCost float32   `json:"OtherCost,omitempty"`
	Projects  []int     `json:"ProjectIds,omitempty"`
	StartDate time.Time `json:"StartDate,omitempty"`
}

type IMemberService interface {
	GetAll(ctx *gin.Context) ([]Member, error)
	CreateMember(payload *CreateMemberRequest, ctx *gin.Context) (Member, error)
	UpdateMember(id int, payload *CreateMemberRequest, ctx *gin.Context) (Member, error)
	GetMemberById(id int) (Member, error)
	DeleteMember(id int) error
	AssignProjectsToMember(memberId int, projectIds []int) error
	AssignRolesToMember(memberId int, roleIds []int) error
}
