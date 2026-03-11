package db

import "gms/backend/internal/models"

func AutoMigrate(dbConn interface {
	AutoMigrate(dst ...interface{}) error
}) error {
	return dbConn.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.Task{},
		&models.TaskTransitionLog{},
		&models.Customer{},
		&models.Lead{},
	)
}
