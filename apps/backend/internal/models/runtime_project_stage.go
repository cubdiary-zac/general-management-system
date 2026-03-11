package models

import "time"

type RuntimeProjectStageStatus string

const (
	RuntimeProjectStageStatusPending RuntimeProjectStageStatus = "pending"
	RuntimeProjectStageStatusActive  RuntimeProjectStageStatus = "active"
	RuntimeProjectStageStatusDone    RuntimeProjectStageStatus = "done"
)

func (s RuntimeProjectStageStatus) IsValid() bool {
	switch s {
	case RuntimeProjectStageStatusPending, RuntimeProjectStageStatusActive, RuntimeProjectStageStatusDone:
		return true
	default:
		return false
	}
}

type RuntimeProjectStage struct {
	ID               uint                      `gorm:"primaryKey" json:"id"`
	CreatedAt        time.Time                 `json:"createdAt"`
	UpdatedAt        time.Time                 `json:"updatedAt"`
	RuntimeProjectID uint                      `gorm:"index;not null" json:"runtimeProjectId"`
	StageTemplateID  uint                      `gorm:"index;not null" json:"stageTemplateId"`
	Name             string                    `gorm:"size:160;not null" json:"name"`
	Code             string                    `gorm:"size:80;not null" json:"code"`
	Description      string                    `gorm:"size:1000" json:"description"`
	Position         int                       `gorm:"not null;default:1" json:"position"`
	Status           RuntimeProjectStageStatus `gorm:"type:varchar(20);not null;default:pending" json:"status"`
}
