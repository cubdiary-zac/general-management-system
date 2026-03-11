package models

import "time"

type Customer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `gorm:"size:120;not null" json:"name"`
	Company   string    `gorm:"size:160" json:"company"`
	Email     string    `gorm:"size:160" json:"email"`
	Phone     string    `gorm:"size:40" json:"phone"`
	OwnerID   uint      `gorm:"index;not null" json:"ownerId"`
}
