package main

import (
	"log"

	"gms/backend/internal/config"
	"gms/backend/internal/db"
	"gms/backend/internal/handlers"
)

func main() {
	cfg := config.Load()

	dbConn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(dbConn); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	if err := db.SeedAdminUser(dbConn, cfg); err != nil {
		log.Fatalf("failed to seed admin user: %v", err)
	}

	router := handlers.SetupRouter(dbConn, cfg)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
