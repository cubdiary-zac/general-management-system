package models

import "time"

type RuntimeProject struct {
	ID                     uint      `gorm:"primaryKey" json:"id"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
	Name                   string    `gorm:"size:160;not null" json:"name"`
	Description            string    `gorm:"size:1000" json:"description"`
	IndustryTemplateID     uint      `gorm:"index;not null" json:"industryTemplateId"`
	ProjectTemplateID      uint      `gorm:"index;not null" json:"projectTemplateId"`
	ProjectTemplateVersion int       `gorm:"not null" json:"projectTemplateVersion"`
	CreatedBy              uint      `gorm:"index;not null" json:"createdBy"`
}
