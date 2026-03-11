package models

import "time"

type RuntimeProjectField struct {
	ID                   uint                `gorm:"primaryKey" json:"id"`
	CreatedAt            time.Time           `json:"createdAt"`
	UpdatedAt            time.Time           `json:"updatedAt"`
	RuntimeProjectID     uint                `gorm:"index;not null" json:"runtimeProjectId"`
	RuntimeProjectFormID uint                `gorm:"index;not null" json:"runtimeProjectFormId"`
	FormFieldTemplateID  uint                `gorm:"index;not null" json:"formFieldTemplateId"`
	Name                 string              `gorm:"size:160;not null" json:"name"`
	Code                 string              `gorm:"size:80;not null" json:"code"`
	Description          string              `gorm:"size:1000" json:"description"`
	Position             int                 `gorm:"not null;default:1" json:"position"`
	WidgetType           FormFieldWidgetType `gorm:"type:varchar(30);not null" json:"widgetType"`
	ValueText            *string             `gorm:"size:4000" json:"valueText"`
}
