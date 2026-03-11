package models

import "time"

type FormFieldWidgetType string

const (
	FormFieldWidgetInput      FormFieldWidgetType = "input"
	FormFieldWidgetTextarea   FormFieldWidgetType = "textarea"
	FormFieldWidgetAttachment FormFieldWidgetType = "attachment"
	FormFieldWidgetSelect     FormFieldWidgetType = "select"
	FormFieldWidgetDate       FormFieldWidgetType = "date"
)

func (w FormFieldWidgetType) IsValid() bool {
	switch w {
	case FormFieldWidgetInput, FormFieldWidgetTextarea, FormFieldWidgetAttachment, FormFieldWidgetSelect, FormFieldWidgetDate:
		return true
	default:
		return false
	}
}

type FormFieldTemplate struct {
	ID             uint                `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time           `json:"createdAt"`
	UpdatedAt      time.Time           `json:"updatedAt"`
	FormTemplateID uint                `gorm:"index;not null" json:"formTemplateId"`
	Name           string              `gorm:"size:160;not null" json:"name"`
	Code           string              `gorm:"size:80;not null" json:"code"`
	Description    string              `gorm:"size:1000" json:"description"`
	Version        int                 `gorm:"not null;default:1" json:"version"`
	Status         TemplateStatus      `gorm:"type:varchar(20);not null;default:draft" json:"status"`
	Position       int                 `gorm:"not null;default:1" json:"position"`
	WidgetType     FormFieldWidgetType `gorm:"type:varchar(30);not null" json:"widgetType"`
}
