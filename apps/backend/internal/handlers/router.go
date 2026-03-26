package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gms/backend/internal/config"
	"gms/backend/internal/middleware"
)

func SetupRouter(db *gorm.DB, cfg config.Config) *gin.Engine {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	authHandler := NewAuthHandler(db, cfg)
	feishuWebhookHandler := NewFeishuWebhookHandler(cfg)

	api := router.Group("/api")
	{
		api.GET("/health", Health)
		api.POST("/feishu/callback", feishuWebhookHandler.Callback)

		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", authHandler.Login)
			authGroup.GET("/me", middleware.AuthRequired(cfg.JWTSecret, db), authHandler.Me)
		}

		registerModules(api, db, cfg)
	}

	return router
}
