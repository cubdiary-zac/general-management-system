package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gms/backend/internal/config"
)

type routeModule interface {
	Name() string
	Register(api *gin.RouterGroup, db *gorm.DB, cfg config.Config)
}

func registerModules(api *gin.RouterGroup, db *gorm.DB, cfg config.Config) {
	modules := []routeModule{
		pmRouteModule{},
		crmRouteModule{},
		stubRouteModule{name: "hr"},
		stubRouteModule{name: "fin"},
	}

	for _, module := range modules {
		module.Register(api, db, cfg)
	}
}
