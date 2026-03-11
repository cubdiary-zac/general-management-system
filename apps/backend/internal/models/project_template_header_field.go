package models

import "time"

type ProjectTemplateHeaderField struct {
	ID                uint                `gorm:"primaryKey" json:"id"`
	CreatedAt         time.Time           `json:"createdAt"`
	UpdatedAt         time.Time           `json:"updatedAt"`
	ProjectTemplateID uint                `gorm:"index;not null" json:"projectTemplateId"`
	Name              string              `gorm:"size:160;not null" json:"name"`
	Code              string              `gorm:"size:80;not null" json:"code"`
	WidgetType        FormFieldWidgetType `gorm:"type:varchar(30);not null" json:"widgetType"`
	Required          bool                `gorm:"not null;default:false" json:"required"`
	Position          int                 `gorm:"not null;default:1" json:"position"`
}
