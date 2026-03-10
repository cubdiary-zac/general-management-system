package models

import "time"

type TaskStatus string

const (
	TaskTodo       TaskStatus = "todo"
	TaskInProgress TaskStatus = "in_progress"
	TaskInReview   TaskStatus = "in_review"
	TaskDone       TaskStatus = "done"
)

type Task struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	ProjectID   uint       `gorm:"index;not null" json:"projectId"`
	Title       string     `gorm:"size:180;not null" json:"title"`
	Description string     `gorm:"size:1000" json:"description"`
	AssigneeID  *uint      `gorm:"index" json:"assigneeId,omitempty"`
	Status      TaskStatus `gorm:"type:varchar(20);not null;default:todo" json:"status"`
}

func (s TaskStatus) IsValid() bool {
	switch s {
	case TaskTodo, TaskInProgress, TaskInReview, TaskDone:
		return true
	default:
		return false
	}
}

func CanTransitionTaskStatus(from TaskStatus, to TaskStatus) bool {
	if !from.IsValid() || !to.IsValid() {
		return false
	}

	if from == to {
		return true
	}

	transitions := map[TaskStatus]TaskStatus{
		TaskTodo:       TaskInProgress,
		TaskInProgress: TaskInReview,
		TaskInReview:   TaskDone,
	}

	next, exists := transitions[from]
	if !exists {
		return false
	}

	return next == to
}
