package models

import "time"

type IndustryTemplate struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	Name        string         `gorm:"size:120;not null" json:"name"`
	Code        string         `gorm:"size:80;not null" json:"code"`
	Description string         `gorm:"size:1000" json:"description"`
	Version     int            `gorm:"not null;default:1" json:"version"`
	Status      TemplateStatus `gorm:"type:varchar(20);not null;default:draft" json:"status"`
}
