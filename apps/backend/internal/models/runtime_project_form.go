package models

import "time"

type RuntimeProjectForm struct {
	ID                    uint      `gorm:"primaryKey" json:"id"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
	RuntimeProjectID      uint      `gorm:"index;not null" json:"runtimeProjectId"`
	RuntimeProjectStageID uint      `gorm:"index;not null" json:"runtimeProjectStageId"`
	FormTemplateID        uint      `gorm:"index;not null" json:"formTemplateId"`
	Name                  string    `gorm:"size:160;not null" json:"name"`
	Code                  string    `gorm:"size:80;not null" json:"code"`
	Description           string    `gorm:"size:1000" json:"description"`
	Position              int       `gorm:"not null;default:1" json:"position"`
}
