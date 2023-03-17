package server

import (
	appctx "http-server/internal/appctx"
	"http-server/internal/controllers"
	"http-server/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	appCtx *appctx.AppCtx
	engine *gin.Engine
}

func NewServer(DB *gorm.DB) *Server {
	args := &appctx.AppCtxArgs{
		DB:                    DB,
		GetProjectServiceFn:   services.GetProjectService,
		GetJwtServiceFn:       services.GetJwtService,
		GetMemberServiceFn:    services.GetMemberService,
		GetTimesheetServiceFn: services.GetTimesheetService,
	}

	return &Server{
		appCtx: appctx.NewAppCtx(args),
		engine: gin.Default(),
	}
}

func (s *Server) Start() {
	s.engine.Use(cors.Default())
	s.RegisterHandler()
	_ = s.engine.Run(":8080")
}

func (s *Server) RegisterHandler() {
	controllers.Register(s.appCtx, s.engine)
}
