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
		ID        int    `json:"id"`
		ProjectID int    `json:"projectId"`
		Status    string `json:"status"`
		Title     string `json:"title"`
	} `json:"items"`
}

type taskDetailResp struct {
	ID        int    `json:"id"`
	ProjectID int    `json:"projectId"`
	Status    string `json:"status"`
	Title     string `json:"title"`
}

type listTaskLogsResp struct {
	Items []struct {
		ID         int    `json:"id"`
		TaskID     int    `json:"taskId"`
		FromStatus string `json:"fromStatus"`
		ToStatus   string `json:"toStatus"`
		OperatorID int    `json:"operatorId"`
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

func TestBlackbox_PMTaskFiltersDetailAndLogs(t *testing.T) {
	h, _, _ := setupTestRouter(t)

	login := doJSONRequest(t, h, http.MethodPost, "/api/auth/login", map[string]any{
		"email":    "admin@gms.local",
		"password": "admin123",
	}, "")
	if login.Code != http.StatusOK {
		t.Fatalf("login failed: %d body=%s", login.Code, login.Body.String())
	}
	loginPayload := decodeJSON[loginResp](t, login)
	token := loginPayload.Token

	projectA := doJSONRequest(t, h, http.MethodPost, "/api/pm/projects", map[string]any{
		"name": "M3 Filters A",
	}, token)
	if projectA.Code != http.StatusCreated {
		t.Fatalf("create project A failed: %d body=%s", projectA.Code, projectA.Body.String())
	}
	projectAID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, projectA).ID

	projectB := doJSONRequest(t, h, http.MethodPost, "/api/pm/projects", map[string]any{
		"name": "M3 Filters B",
	}, token)
	if projectB.Code != http.StatusCreated {
		t.Fatalf("create project B failed: %d body=%s", projectB.Code, projectB.Body.String())
	}
	projectBID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, projectB).ID

	taskA1 := doJSONRequest(t, h, http.MethodPost, "/api/pm/tasks", map[string]any{
		"projectId": projectAID,
		"title":     "alpha wireframes",
	}, token)
	if taskA1.Code != http.StatusCreated {
		t.Fatalf("create task A1 failed: %d body=%s", taskA1.Code, taskA1.Body.String())
	}
	taskA1ID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, taskA1).ID

	taskA2 := doJSONRequest(t, h, http.MethodPost, "/api/pm/tasks", map[string]any{
		"projectId": projectAID,
		"title":     "beta docs",
	}, token)
	if taskA2.Code != http.StatusCreated {
		t.Fatalf("create task A2 failed: %d body=%s", taskA2.Code, taskA2.Body.String())
	}

	taskB1 := doJSONRequest(t, h, http.MethodPost, "/api/pm/tasks", map[string]any{
		"projectId": projectBID,
		"title":     "alpha follow-up",
	}, token)
	if taskB1.Code != http.StatusCreated {
		t.Fatalf("create task B1 failed: %d body=%s", taskB1.Code, taskB1.Body.String())
	}

	patch := doJSONRequest(t, h, http.MethodPatch, "/api/pm/tasks/"+itoa(taskA1ID)+"/status", map[string]any{
		"status": "in_progress",
	}, token)
	if patch.Code != http.StatusOK {
		t.Fatalf("patch status failed: %d body=%s", patch.Code, patch.Body.String())
	}

	filtered := doJSONRequest(t, h, http.MethodGet, "/api/pm/tasks?projectId="+itoa(projectAID)+"&status=in_progress&q=ALPHA", nil, token)
	if filtered.Code != http.StatusOK {
		t.Fatalf("filtered list failed: %d body=%s", filtered.Code, filtered.Body.String())
	}
	filteredPayload := decodeJSON[listTasksResp](t, filtered)
	if len(filteredPayload.Items) != 1 {
		t.Fatalf("expected 1 filtered task, got %d", len(filteredPayload.Items))
	}
	if filteredPayload.Items[0].ID != taskA1ID {
		t.Fatalf("expected filtered task id %d, got %d", taskA1ID, filteredPayload.Items[0].ID)
	}

	detail := doJSONRequest(t, h, http.MethodGet, "/api/pm/tasks/"+itoa(taskA1ID), nil, token)
	if detail.Code != http.StatusOK {
		t.Fatalf("detail failed: %d body=%s", detail.Code, detail.Body.String())
	}
	detailPayload := decodeJSON[taskDetailResp](t, detail)
	if detailPayload.ID != taskA1ID || detailPayload.Status != "in_progress" {
		t.Fatalf("unexpected detail payload: %+v", detailPayload)
	}

	logs := doJSONRequest(t, h, http.MethodGet, "/api/pm/tasks/"+itoa(taskA1ID)+"/logs", nil, token)
	if logs.Code != http.StatusOK {
		t.Fatalf("logs failed: %d body=%s", logs.Code, logs.Body.String())
	}
	logPayload := decodeJSON[listTaskLogsResp](t, logs)
	if len(logPayload.Items) == 0 {
		t.Fatalf("expected at least one transition log")
	}
	firstLog := logPayload.Items[0]
	if firstLog.TaskID != taskA1ID || firstLog.FromStatus != "todo" || firstLog.ToStatus != "in_progress" {
		t.Fatalf("unexpected transition log payload: %+v", firstLog)
	}
	if firstLog.OperatorID != loginPayload.User.ID {
		t.Fatalf("expected operator %d, got %d", loginPayload.User.ID, firstLog.OperatorID)
	}
}

func TestBlackbox_ModuleStubHealthEndpoints(t *testing.T) {
	h, _, _ := setupTestRouter(t)

	login := doJSONRequest(t, h, http.MethodPost, "/api/auth/login", map[string]any{
		"email":    "admin@gms.local",
		"password": "admin123",
	}, "")
	if login.Code != http.StatusOK {
		t.Fatalf("login failed: %d body=%s", login.Code, login.Body.String())
	}
	token := decodeJSON[loginResp](t, login).Token

	for _, moduleName := range []string{"crm", "hr", "fin"} {
		resp := doJSONRequest(t, h, http.MethodGet, "/api/"+moduleName+"/health", nil, token)
		if resp.Code != http.StatusOK {
			t.Fatalf("module %s health failed: %d body=%s", moduleName, resp.Code, resp.Body.String())
		}
	}
}
