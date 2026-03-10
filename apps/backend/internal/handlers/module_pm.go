package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gms/backend/internal/config"
	"gms/backend/internal/middleware"
	"gms/backend/internal/models"
)

type pmRouteModule struct{}

func (pmRouteModule) Name() string {
	return "pm"
}

func (pmRouteModule) Register(api *gin.RouterGroup, db *gorm.DB, cfg config.Config) {
	pmHandler := NewPMHandler(db)

	pm := api.Group("/pm", middleware.AuthRequired(cfg.JWTSecret, db))
	{
		pm.GET("/projects", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember, models.RoleViewer), pmHandler.ListProjects)
		pm.POST("/projects", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), pmHandler.CreateProject)
		pm.GET("/tasks", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember, models.RoleViewer), pmHandler.ListTasks)
		pm.GET("/tasks/:id", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember, models.RoleViewer), pmHandler.GetTask)
		pm.GET("/tasks/:id/logs", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember, models.RoleViewer), pmHandler.ListTaskLogs)
		pm.POST("/tasks", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), pmHandler.CreateTask)
		pm.PATCH("/tasks/:id/status", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), pmHandler.PatchTaskStatus)
	}
}
