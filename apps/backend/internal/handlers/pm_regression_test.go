package handlers

import (
	"net/http"
	"testing"
)

// Regression: invalid transition should not mutate task status.
func TestRegression_InvalidTransitionDoesNotMutateStatus(t *testing.T) {
	h, _, _ := setupTestRouter(t)

	login := doJSONRequest(t, h, http.MethodPost, "/api/auth/login", map[string]any{
		"email":    "admin@gms.local",
		"password": "admin123",
	}, "")
	if login.Code != http.StatusOK {
		t.Fatalf("login failed: %d", login.Code)
	}
	token := decodeJSON[loginResp](t, login).Token

	project := doJSONRequest(t, h, http.MethodPost, "/api/pm/projects", map[string]any{
		"name": "Regression PM",
	}, token)
	if project.Code != http.StatusCreated {
		t.Fatalf("create project failed: %d body=%s", project.Code, project.Body.String())
	}
	projectID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, project).ID

	task := doJSONRequest(t, h, http.MethodPost, "/api/pm/tasks", map[string]any{
		"projectId": projectID,
		"title":     "must-stay-todo",
	}, token)
	if task.Code != http.StatusCreated {
		t.Fatalf("create task failed: %d body=%s", task.Code, task.Body.String())
	}
	taskID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, task).ID

	invalid := doJSONRequest(t, h, http.MethodPatch, "/api/pm/tasks/"+itoa(taskID)+"/status", map[string]any{
		"status": "done",
	}, token)
	if invalid.Code != http.StatusBadRequest {
		t.Fatalf("expected invalid transition 400, got %d body=%s", invalid.Code, invalid.Body.String())
	}

	list := doJSONRequest(t, h, http.MethodGet, "/api/pm/tasks?projectId="+itoa(projectID), nil, token)
	if list.Code != http.StatusOK {
		t.Fatalf("list tasks failed: %d", list.Code)
	}
	items := decodeJSON[listTasksResp](t, list).Items
	if len(items) == 0 {
		t.Fatalf("expected tasks in list")
	}
	if items[0].Status != "todo" {
		t.Fatalf("expected status todo after invalid transition, got %s", items[0].Status)
	}
}
