package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gms/backend/internal/config"
	"gms/backend/internal/middleware"
	"gms/backend/internal/models"
)

type stubRouteModule struct {
	name string
}

func (m stubRouteModule) Name() string {
	return m.name
}

func (m stubRouteModule) Register(api *gin.RouterGroup, db *gorm.DB, cfg config.Config) {
	group := api.Group("/"+m.name, middleware.AuthRequired(cfg.JWTSecret, db))
	group.GET("/health", middleware.RequireRoles(models.RoleOwner, models.RoleAdmin, models.RoleMember, models.RoleViewer), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"module":  strings.ToUpper(m.name),
			"status":  "ok",
			"message": "placeholder endpoint",
		})
	})
}
