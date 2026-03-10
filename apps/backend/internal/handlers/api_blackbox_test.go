package handlers

import (
	"net/http"
	"strconv"
	"testing"
)

type loginResp struct {
	Token string `json:"token"`
	User  struct {
		ID int `json:"id"`
	} `json:"user"`
}

type listProjectsResp struct {
	Items []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"items"`
}

type listTasksResp struct {
	Items []struct {
		ID       int    `json:"id"`
		ProjectID int   `json:"projectId"`
		Status   string `json:"status"`
		Title    string `json:"title"`
	} `json:"items"`
}

func TestBlackbox_HealthAndAuth(t *testing.T) {
	h, _, _ := setupTestRouter(t)

	rr := doJSONRequest(t, h, http.MethodGet, "/api/health", nil, "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected health 200, got %d", rr.Code)
	}

	login := doJSONRequest(t, h, http.MethodPost, "/api/auth/login", map[string]any{
		"email":    "admin@gms.local",
		"password": "admin123",
	}, "")
	if login.Code != http.StatusOK {
		t.Fatalf("expected login 200, got %d body=%s", login.Code, login.Body.String())
	}

	payload := decodeJSON[loginResp](t, login)
	if payload.Token == "" || payload.User.ID == 0 {
		t.Fatalf("expected token and user in login response")
	}

	me := doJSONRequest(t, h, http.MethodGet, "/api/auth/me", nil, payload.Token)
	if me.Code != http.StatusOK {
		t.Fatalf("expected me 200, got %d body=%s", me.Code, me.Body.String())
	}
}

func itoa(v int) string { return strconv.Itoa(v) }

func TestBlackbox_PMProjectTaskFlow(t *testing.T) {
	h, _, _ := setupTestRouter(t)

	login := doJSONRequest(t, h, http.MethodPost, "/api/auth/login", map[string]any{
		"email":    "admin@gms.local",
		"password": "admin123",
	}, "")
	if login.Code != http.StatusOK {
		t.Fatalf("login failed: %d body=%s", login.Code, login.Body.String())
	}
	token := decodeJSON[loginResp](t, login).Token

	createProject := doJSONRequest(t, h, http.MethodPost, "/api/pm/projects", map[string]any{
		"name":        "M2 Blackbox",
		"description": "project from blackbox test",
	}, token)
	if createProject.Code != http.StatusCreated {
		t.Fatalf("expected create project 201, got %d body=%s", createProject.Code, createProject.Body.String())
	}
	createdProject := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createProject)

	projects := doJSONRequest(t, h, http.MethodGet, "/api/pm/projects", nil, token)
	if projects.Code != http.StatusOK {
		t.Fatalf("expected list projects 200, got %d", projects.Code)
	}
	projPayload := decodeJSON[listProjectsResp](t, projects)
	if len(projPayload.Items) == 0 {
		t.Fatalf("expected at least one project")
	}

	createTask := doJSONRequest(t, h, http.MethodPost, "/api/pm/tasks", map[string]any{
		"projectId":   createdProject.ID,
		"title":       "blackbox-task",
		"description": "task from test",
	}, token)
	if createTask.Code != http.StatusCreated {
		t.Fatalf("expected create task 201, got %d body=%s", createTask.Code, createTask.Body.String())
	}
	createdTask := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createTask)

	tasks := doJSONRequest(t, h, http.MethodGet, "/api/pm/tasks?projectId="+itoa(createdProject.ID), nil, token)
	if tasks.Code != http.StatusOK {
		t.Fatalf("expected list tasks 200, got %d body=%s", tasks.Code, tasks.Body.String())
	}
	taskPayload := decodeJSON[listTasksResp](t, tasks)
	if len(taskPayload.Items) == 0 {
		t.Fatalf("expected at least one task")
	}

	patch := doJSONRequest(t, h, http.MethodPatch, "/api/pm/tasks/"+itoa(createdTask.ID)+"/status", map[string]any{
		"status": "in_progress",
	}, token)
	if patch.Code != http.StatusOK {
		t.Fatalf("expected patch status 200, got %d body=%s", patch.Code, patch.Body.String())
	}
}
