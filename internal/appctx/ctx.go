package appctx

import (
	"http-server/cmd/server/config"
	"http-server/internal/services"
	"http-server/internal/services/auth"
	"http-server/internal/services/member"
	"http-server/internal/services/project"
	"http-server/internal/services/timesheet"

	"gorm.io/gorm"
)

type IAppCtx interface {
	GetProjectService() project.IProjectService
	GetJwtService() auth.IJWT
	GetMemberService() member.IMemberService
	GetTimesheetService() timesheet.ITimeSheet
}

type AppCtxArgs struct {
	DB *gorm.DB

	GetProjectServiceFn   services.GetProjectServiceFn
	GetJwtServiceFn       services.GetJwtServiceFn
	GetMemberServiceFn    services.GetMemberServiceFn
	GetTimesheetServiceFn services.GetTimesheetServiceFn
}

type AppCtx struct {
	DB *gorm.DB

	projectService   project.IProjectService
	jwtService       auth.IJWT
	memberService    member.IMemberService
	timesheetService timesheet.ITimeSheet
}

func NewAppCtx(args *AppCtxArgs) *AppCtx {
	config, _ := config.LoadConfig(".")

	return &AppCtx{
		DB:             args.DB,
		projectService: args.GetProjectServiceFn(args.DB),
		jwtService:     args.GetJwtServiceFn(&config),
		memberService:  args.GetMemberServiceFn(args.DB),
		timesheetService: args.GetTimesheetServiceFn(args.DB),
	}
}

func (a *AppCtx) GetProjectService() project.IProjectService {
	return a.projectService
}

func (a *AppCtx) GetJwtService() auth.IJWT {
	return a.jwtService
}

func (a *AppCtx) GetMemberService() member.IMemberService {
	return a.memberService
}

func (a *AppCtx) GetTimesheetService() timesheet.ITimeSheet {
	return a.timesheetService
}

func (a *AppCtx) GetUserService() {}
