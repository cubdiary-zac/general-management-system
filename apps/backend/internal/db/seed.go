package db

import (
	"errors"

	"gorm.io/gorm"

	"gms/backend/internal/auth"
	"gms/backend/internal/config"
	"gms/backend/internal/models"
)

func SeedAdminUser(dbConn *gorm.DB, cfg config.Config) error {
	var existing models.User
	err := dbConn.Where("email = ?", cfg.SeedAdminEmail).First(&existing).Error
	if err == nil {
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hash, err := auth.HashPassword(cfg.SeedAdminPassword)
	if err != nil {
		return err
	}

	admin := models.User{
		Name:         cfg.SeedAdminName,
		Email:        cfg.SeedAdminEmail,
		PasswordHash: hash,
		Role:         models.RoleOwner,
	}

	return dbConn.Create(&admin).Error
}
