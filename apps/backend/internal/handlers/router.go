package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gms/backend/internal/config"
	"gms/backend/internal/middleware"
	"gms/backend/internal/models"
)

func SetupRouter(db *gorm.DB, cfg config.Config) *gin.Engine {
	router := gin.Default()

	authHandler := NewAuthHandler(db, cfg)
	pmHandler := NewPMHandler(db)

	api := router.Group("/api")
	{
		api.GET("/health", Health)

		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", authHandler.Login)
			authGroup.GET("/me", middleware.AuthRequired(cfg.JWTSecret, db), authHandler.Me)
		}

		pm := api.Group("/pm", middleware.AuthRequired(cfg.JWTSecret, db))
		{
			pm.GET("/projects", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember, models.RoleViewer), pmHandler.ListProjects)
			pm.POST("/projects", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), pmHandler.CreateProject)
			pm.GET("/tasks", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember, models.RoleViewer), pmHandler.ListTasks)
			pm.POST("/tasks", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), pmHandler.CreateTask)
			pm.PATCH("/tasks/:id/status", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), pmHandler.PatchTaskStatus)
		}
	}

	return router
}
