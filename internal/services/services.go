package services

import (
	"http-server/cmd/server/config"
	"http-server/internal/services/auth"
	"http-server/internal/services/member"
	"http-server/internal/services/project"
	"http-server/internal/services/timesheet"

	"gorm.io/gorm"
)

type GetProjectServiceFn func(DB *gorm.DB) project.IProjectService
type GetJwtServiceFn func(config *config.Config) auth.IJWT
type GetMemberServiceFn func(*gorm.DB) member.IMemberService
type GetTimesheetServiceFn func(*gorm.DB) timesheet.ITimeSheet

func GetProjectService(DB *gorm.DB) project.IProjectService {
	return project.NewProjectService(DB)
}

func GetJwtService(config *config.Config) auth.IJWT {
	return auth.NewJWT(config)
}

func GetMemberService(DB *gorm.DB) member.IMemberService {
	return member.NewMemberService(DB)
}

func GetTimesheetService(DB *gorm.DB) timesheet.ITimeSheet {
	return timesheet.NewTimeSheetService(DB)
}
