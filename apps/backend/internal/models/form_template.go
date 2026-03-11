package models

import "time"

type FormTemplate struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	StageTemplateID uint           `gorm:"index;not null" json:"stageTemplateId"`
	Name            string         `gorm:"size:160;not null" json:"name"`
	Code            string         `gorm:"size:80;not null" json:"code"`
	Description     string         `gorm:"size:1000" json:"description"`
	Version         int            `gorm:"not null;default:1" json:"version"`
	Status          TemplateStatus `gorm:"type:varchar(20);not null;default:draft" json:"status"`
	PublishedAt     *time.Time     `json:"publishedAt,omitempty"`
	PublishedBy     *uint          `gorm:"index" json:"publishedBy,omitempty"`
	Position        int            `gorm:"not null;default:1" json:"position"`
}
