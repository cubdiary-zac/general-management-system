package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gms/backend/internal/config"
	"gms/backend/internal/middleware"
	"gms/backend/internal/models"
)

type crmRouteModule struct{}

func (crmRouteModule) Name() string {
	return "crm"
}

func (crmRouteModule) Register(api *gin.RouterGroup, db *gorm.DB, cfg config.Config) {
	crmHandler := NewCRMHandler(db)

	crm := api.Group("/crm", middleware.AuthRequired(cfg.JWTSecret, db))
	{
		crm.GET("/customers", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember, models.RoleViewer), crmHandler.ListCustomers)
		crm.POST("/customers", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), crmHandler.CreateCustomer)
		crm.GET("/leads", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember, models.RoleViewer), crmHandler.ListLeads)
		crm.POST("/leads", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), crmHandler.CreateLead)
		crm.PATCH("/leads/:id/status", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember), crmHandler.PatchLeadStatus)
		crm.GET("/summary", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember, models.RoleViewer), crmHandler.Summary)
	}
}
