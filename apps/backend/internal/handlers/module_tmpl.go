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
		tmpl.POST("/industries/:id/publish", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.PublishIndustryTemplate)
		tmpl.POST("/industries/:id/unpublish", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.UnpublishIndustryTemplate)
		tmpl.GET("/project-templates", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember, models.RoleViewer), tmplHandler.ListProjectTemplates)
		tmpl.POST("/project-templates", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.CreateProjectTemplate)
		tmpl.GET("/project-templates/:id/header-fields", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember, models.RoleViewer), tmplHandler.ListProjectTemplateHeaderFields)
		tmpl.POST("/project-templates/:id/header-fields", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.CreateProjectTemplateHeaderField)
		tmpl.PATCH("/project-templates/:id/header-fields/:fieldId", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.PatchProjectTemplateHeaderField)
		tmpl.POST("/project-templates/:id/publish", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.PublishProjectTemplate)
		tmpl.POST("/project-templates/:id/unpublish", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.UnpublishProjectTemplate)
		tmpl.POST("/project-templates/:id/next-version", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.CreateNextProjectTemplateVersion)
		tmpl.GET("/project-templates/:id/lifecycle", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember, models.RoleViewer), tmplHandler.GetProjectTemplateLifecycle)
		tmpl.POST("/project-templates/:id/instantiate", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.InstantiateProjectTemplate)
		tmpl.GET("/stage-templates", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember, models.RoleViewer), tmplHandler.ListStageTemplates)
		tmpl.POST("/stage-templates", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.CreateStageTemplate)
		tmpl.PATCH("/stage-templates/:id", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.PatchStageTemplate)
		tmpl.POST("/stage-templates/:id/publish", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.PublishStageTemplate)
		tmpl.POST("/stage-templates/:id/unpublish", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.UnpublishStageTemplate)
		tmpl.GET("/form-templates", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember, models.RoleViewer), tmplHandler.ListFormTemplates)
		tmpl.POST("/form-templates", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.CreateFormTemplate)
		tmpl.POST("/form-templates/:id/publish", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.PublishFormTemplate)
		tmpl.POST("/form-templates/:id/unpublish", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.UnpublishFormTemplate)
		tmpl.GET("/field-templates", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember, models.RoleViewer), tmplHandler.ListFormFieldTemplates)
		tmpl.POST("/field-templates", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.CreateFormFieldTemplate)
		tmpl.POST("/field-templates/:id/publish", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.PublishFormFieldTemplate)
		tmpl.POST("/field-templates/:id/unpublish", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), tmplHandler.UnpublishFormFieldTemplate)
	}
}
