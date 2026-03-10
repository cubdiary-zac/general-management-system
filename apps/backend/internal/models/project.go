package models

import "time"

type Project struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Name        string    `gorm:"size:120;not null" json:"name"`
	Description string    `gorm:"size:500" json:"description"`
	OwnerID     uint      `gorm:"not null" json:"ownerId"`
}
