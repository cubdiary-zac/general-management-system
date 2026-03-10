package handlers

import (
	"net/http"
	"testing"
)

// Smoke: shortest happy path chain proving system is runnable.
func TestSmoke_HappyPathChain(t *testing.T) {
	h, _, _ := setupTestRouter(t)

	health := doJSONRequest(t, h, http.MethodGet, "/api/health", nil, "")
	if health.Code != http.StatusOK {
		t.Fatalf("health failed: %d", health.Code)
	}

	login := doJSONRequest(t, h, http.MethodPost, "/api/auth/login", map[string]any{
		"email":    "admin@gms.local",
		"password": "admin123",
	}, "")
	if login.Code != http.StatusOK {
		t.Fatalf("login failed: %d body=%s", login.Code, login.Body.String())
	}
	token := decodeJSON[loginResp](t, login).Token

	project := doJSONRequest(t, h, http.MethodPost, "/api/pm/projects", map[string]any{"name": "Smoke"}, token)
	if project.Code != http.StatusCreated {
		t.Fatalf("create project failed: %d", project.Code)
	}
	projectID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, project).ID

	task := doJSONRequest(t, h, http.MethodPost, "/api/pm/tasks", map[string]any{
		"projectId": projectID,
		"title":     "smoke-task",
	}, token)
	if task.Code != http.StatusCreated {
		t.Fatalf("create task failed: %d", task.Code)
	}
	taskID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, task).ID

	steps := []string{"in_progress", "in_review", "done"}
	for _, step := range steps {
		r := doJSONRequest(t, h, http.MethodPatch, "/api/pm/tasks/"+itoa(taskID)+"/status", map[string]any{
			"status": step,
		}, token)
		if r.Code != http.StatusOK {
			t.Fatalf("step %s failed: %d body=%s", step, r.Code, r.Body.String())
		}
	}
}
