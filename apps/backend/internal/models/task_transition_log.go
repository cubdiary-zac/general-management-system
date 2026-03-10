package models

import "time"

type TaskTransitionLog struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time  `json:"createdAt"`
	TaskID     uint       `gorm:"index;not null" json:"taskId"`
	FromStatus TaskStatus `gorm:"type:varchar(20);not null" json:"fromStatus"`
	ToStatus   TaskStatus `gorm:"type:varchar(20);not null" json:"toStatus"`
	OperatorID uint       `gorm:"index;not null" json:"operatorId"`
}
