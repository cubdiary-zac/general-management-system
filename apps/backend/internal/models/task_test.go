package models

import "testing"

func TestCanTransitionTaskStatus(t *testing.T) {
	tests := []struct {
		name     string
		from     TaskStatus
		to       TaskStatus
		expected bool
	}{
		{name: "todo to in_progress", from: TaskTodo, to: TaskInProgress, expected: true},
		{name: "in_progress to in_review", from: TaskInProgress, to: TaskInReview, expected: true},
		{name: "in_review to done", from: TaskInReview, to: TaskDone, expected: true},
		{name: "same status", from: TaskDone, to: TaskDone, expected: true},
		{name: "todo to done rejected", from: TaskTodo, to: TaskDone, expected: false},
		{name: "done to in_progress rejected", from: TaskDone, to: TaskInProgress, expected: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := CanTransitionTaskStatus(tc.from, tc.to)
			if result != tc.expected {
				t.Fatalf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}
