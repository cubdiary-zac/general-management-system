package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gms/backend/internal/config"
	"gms/backend/internal/middleware"
	"gms/backend/internal/models"
)

type tmplRouteModule struct{}

func (tmplRouteModule) Name() string {
	return "tmpl"
}

func (tmplRouteModule) Register(api *gin.RouterGroup, db *gorm.DB, cfg config.Config) {
	tmplHandler := NewTemplateHandler(db)

	tmpl := api.Group("/tmpl", middleware.AuthRequired(cfg.JWTSecret, db))
	{
		tmpl.GET("/industries", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember, models.RoleViewer), tmplHandler.ListIndustryTemplates)
		tmpl.POST("/industries", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.CreateIndustryTemplate)
		tmpl.GET("/project-templates", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember, models.RoleViewer), tmplHandler.ListProjectTemplates)
		tmpl.POST("/project-templates", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.CreateProjectTemplate)
		tmpl.GET("/stage-templates", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember, models.RoleViewer), tmplHandler.ListStageTemplates)
		tmpl.POST("/stage-templates", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.CreateStageTemplate)
		tmpl.GET("/form-templates", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember, models.RoleViewer), tmplHandler.ListFormTemplates)
		tmpl.POST("/form-templates", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.CreateFormTemplate)
		tmpl.GET("/field-templates", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember, models.RoleViewer), tmplHandler.ListFormFieldTemplates)
		tmpl.POST("/field-templates", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.CreateFormFieldTemplate)
	}
}
