package models

import "time"

type LeadStatus string

const (
	LeadNew       LeadStatus = "new"
	LeadContacted LeadStatus = "contacted"
	LeadQualified LeadStatus = "qualified"
	LeadWon       LeadStatus = "won"
	LeadLost      LeadStatus = "lost"
)

type Lead struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	CustomerID *uint      `gorm:"index" json:"customerId,omitempty"`
	Name       string     `gorm:"size:160;not null" json:"name"`
	Source     string     `gorm:"size:80;not null" json:"source"`
	Status     LeadStatus `gorm:"type:varchar(20);not null;default:new" json:"status"`
	Amount     *float64   `json:"amount,omitempty"`
	OwnerID    uint       `gorm:"index;not null" json:"ownerId"`
}

func (s LeadStatus) IsValid() bool {
	switch s {
	case LeadNew, LeadContacted, LeadQualified, LeadWon, LeadLost:
		return true
	default:
		return false
	}
}

func CanTransitionLeadStatus(from LeadStatus, to LeadStatus) bool {
	if !from.IsValid() || !to.IsValid() {
		return false
	}

	if from == to {
		return true
	}

	transitions := map[LeadStatus][]LeadStatus{
		LeadNew:       {LeadContacted, LeadLost},
		LeadContacted: {LeadQualified, LeadLost},
		LeadQualified: {LeadWon, LeadLost},
		LeadWon:       {},
		LeadLost:      {},
	}

	next, exists := transitions[from]
	if !exists {
		return false
	}

	for _, status := range next {
		if status == to {
			return true
		}
	}

	return false
}
