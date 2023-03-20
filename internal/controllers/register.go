package controllers

import (
	"fmt"
	"http-server/internal/appctx"
	"http-server/internal/controllers/auth"
	"http-server/internal/controllers/member"
	"http-server/internal/controllers/project"
	"http-server/internal/controllers/timesheet"

	"github.com/gin-gonic/gin"
)

func Register(appCtx *appctx.AppCtx, engine *gin.Engine) {
	routes := engine.Group("/api/v1")
	routes.GET("/healthcheck", func(ctx *gin.Context) {
		fmt.Print("Status OK")
		ctx.JSON(200, gin.H{"St": "Nice"})
	})

	authHandler := auth.NewAuthHandler(appCtx)
	routes.POST("/login", authHandler.Login)
	routes.POST("/signup", authHandler.Signup)

	// routes.Use(middleware.MiddleWare())

	projectHandler := project.NewProjectHandler(appCtx)
	projectRoutes := routes.Group("/projects")
	projectRoutes.GET("", projectHandler.RetrieveProject)
	projectRoutes.GET("/:id", projectHandler.GetProjectByIdOrName)
	projectRoutes.POST("", projectHandler.CreateProject)
	projectRoutes.PUT("/:id", projectHandler.UpdateProject)
	projectRoutes.DELETE("/:id", projectHandler.DeleteProject)

	memberHandler := member.NewMemberHandler(appCtx)
	memberRoutes := routes.Group("/members")
	memberRoutes.GET("", memberHandler.RetreiveMembers)
	memberRoutes.GET("/:id", memberHandler.GetMember)
	memberRoutes.POST("", memberHandler.CreateMember)
	memberRoutes.PUT("/:id", memberHandler.UpdateMember)
	memberRoutes.DELETE("/:id", memberHandler.DeleteMember)

	timesheetHandler := timesheet.NewTimesheetHandler(appCtx)
	timesheetRoutes := routes.Group("/timesheet")
	timesheetRoutes.GET("", timesheetHandler.GetAll)
	// timesheetRoutes.GET("/:id", timesheetHandler.GetByMemberId)
	timesheetRoutes.GET("/:id", timesheetHandler.GetById)
	timesheetRoutes.POST("", timesheetHandler.Create)
	timesheetRoutes.POST("/:id/create_segment", timesheetHandler.SaveSegment)
	timesheetRoutes.POST("/create_multiple", timesheetHandler.SaveMultipleTsSegment)
	timesheetRoutes.DELETE("/:id", timesheetHandler.Delete)
}
