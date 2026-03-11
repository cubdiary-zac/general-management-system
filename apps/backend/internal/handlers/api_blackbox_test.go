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

type listCustomersResp struct {
	Items []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"items"`
}

type listLeadsResp struct {
	Items []struct {
		ID         int    `json:"id"`
		CustomerID *int   `json:"customerId"`
		Name       string `json:"name"`
		Status     string `json:"status"`
	} `json:"items"`
}

type crmSummaryResp struct {
	Counts map[string]int `json:"counts"`
}

type listIndustryTemplatesResp struct {
	Items []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"items"`
}

type listProjectTemplatesResp struct {
	Items []struct {
		ID                 int `json:"id"`
		IndustryTemplateID int `json:"industryTemplateId"`
	} `json:"items"`
}

type listStageTemplatesResp struct {
	Items []struct {
		ID                int `json:"id"`
		ProjectTemplateID int `json:"projectTemplateId"`
		Position          int `json:"position"`
	} `json:"items"`
}

type listFormTemplatesResp struct {
	Items []struct {
		ID              int `json:"id"`
		StageTemplateID int `json:"stageTemplateId"`
		Position        int `json:"position"`
	} `json:"items"`
}

type listFieldTemplatesResp struct {
	Items []struct {
		ID             int    `json:"id"`
		FormTemplateID int    `json:"formTemplateId"`
		Position       int    `json:"position"`
		WidgetType     string `json:"widgetType"`
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

func TestBlackbox_CRMCustomerLeadFlow(t *testing.T) {
	h, _, _ := setupTestRouter(t)

	login := doJSONRequest(t, h, http.MethodPost, "/api/auth/login", map[string]any{
		"email":    "admin@gms.local",
		"password": "admin123",
	}, "")
	if login.Code != http.StatusOK {
		t.Fatalf("login failed: %d body=%s", login.Code, login.Body.String())
	}
	token := decodeJSON[loginResp](t, login).Token

	createCustomer := doJSONRequest(t, h, http.MethodPost, "/api/crm/customers", map[string]any{
		"name":    "Acme Ops",
		"company": "Acme Ltd",
		"email":   "ops@acme.local",
		"phone":   "+1-555-1000",
	}, token)
	if createCustomer.Code != http.StatusCreated {
		t.Fatalf("create customer failed: %d body=%s", createCustomer.Code, createCustomer.Body.String())
	}
	customerID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createCustomer).ID

	customers := doJSONRequest(t, h, http.MethodGet, "/api/crm/customers", nil, token)
	if customers.Code != http.StatusOK {
		t.Fatalf("list customers failed: %d body=%s", customers.Code, customers.Body.String())
	}
	customerItems := decodeJSON[listCustomersResp](t, customers).Items
	if len(customerItems) == 0 {
		t.Fatalf("expected at least one customer")
	}

	leadOne := doJSONRequest(t, h, http.MethodPost, "/api/crm/leads", map[string]any{
		"name":       "North Alpha",
		"source":     "referral",
		"customerId": customerID,
		"amount":     18000,
	}, token)
	if leadOne.Code != http.StatusCreated {
		t.Fatalf("create lead one failed: %d body=%s", leadOne.Code, leadOne.Body.String())
	}
	leadOneID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, leadOne).ID

	leadTwo := doJSONRequest(t, h, http.MethodPost, "/api/crm/leads", map[string]any{
		"name":   "South Beta",
		"source": "web",
	}, token)
	if leadTwo.Code != http.StatusCreated {
		t.Fatalf("create lead two failed: %d body=%s", leadTwo.Code, leadTwo.Body.String())
	}
	leadTwoID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, leadTwo).ID

	patchLeadOneContacted := doJSONRequest(t, h, http.MethodPatch, "/api/crm/leads/"+itoa(leadOneID)+"/status", map[string]any{
		"status": "contacted",
	}, token)
	if patchLeadOneContacted.Code != http.StatusOK {
		t.Fatalf("patch lead one contacted failed: %d body=%s", patchLeadOneContacted.Code, patchLeadOneContacted.Body.String())
	}

	patchLeadOneQualified := doJSONRequest(t, h, http.MethodPatch, "/api/crm/leads/"+itoa(leadOneID)+"/status", map[string]any{
		"status": "qualified",
	}, token)
	if patchLeadOneQualified.Code != http.StatusOK {
		t.Fatalf("patch lead one qualified failed: %d body=%s", patchLeadOneQualified.Code, patchLeadOneQualified.Body.String())
	}

	patchLeadOneSame := doJSONRequest(t, h, http.MethodPatch, "/api/crm/leads/"+itoa(leadOneID)+"/status", map[string]any{
		"status": "qualified",
	}, token)
	if patchLeadOneSame.Code != http.StatusOK {
		t.Fatalf("patch lead one same status failed: %d body=%s", patchLeadOneSame.Code, patchLeadOneSame.Body.String())
	}

	patchLeadTwoLost := doJSONRequest(t, h, http.MethodPatch, "/api/crm/leads/"+itoa(leadTwoID)+"/status", map[string]any{
		"status": "lost",
	}, token)
	if patchLeadTwoLost.Code != http.StatusOK {
		t.Fatalf("patch lead two lost failed: %d body=%s", patchLeadTwoLost.Code, patchLeadTwoLost.Body.String())
	}

	filteredLeads := doJSONRequest(t, h, http.MethodGet, "/api/crm/leads?status=qualified&q=ALPHA", nil, token)
	if filteredLeads.Code != http.StatusOK {
		t.Fatalf("filtered leads failed: %d body=%s", filteredLeads.Code, filteredLeads.Body.String())
	}
	filteredItems := decodeJSON[listLeadsResp](t, filteredLeads).Items
	if len(filteredItems) != 1 {
		t.Fatalf("expected 1 filtered lead, got %d", len(filteredItems))
	}
	if filteredItems[0].ID != leadOneID || filteredItems[0].Status != "qualified" {
		t.Fatalf("unexpected filtered lead payload: %+v", filteredItems[0])
	}

	summary := doJSONRequest(t, h, http.MethodGet, "/api/crm/summary", nil, token)
	if summary.Code != http.StatusOK {
		t.Fatalf("summary failed: %d body=%s", summary.Code, summary.Body.String())
	}
	summaryPayload := decodeJSON[crmSummaryResp](t, summary)
	if summaryPayload.Counts["qualified"] != 1 {
		t.Fatalf("expected qualified count 1, got %d", summaryPayload.Counts["qualified"])
	}
	if summaryPayload.Counts["lost"] != 1 {
		t.Fatalf("expected lost count 1, got %d", summaryPayload.Counts["lost"])
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

	for _, moduleName := range []string{"hr", "fin"} {
		resp := doJSONRequest(t, h, http.MethodGet, "/api/"+moduleName+"/health", nil, token)
		if resp.Code != http.StatusOK {
			t.Fatalf("module %s health failed: %d body=%s", moduleName, resp.Code, resp.Body.String())
		}
	}
}

func TestBlackbox_TemplateEngineCreateListFlow(t *testing.T) {
	h, _, _ := setupTestRouter(t)

	login := doJSONRequest(t, h, http.MethodPost, "/api/auth/login", map[string]any{
		"email":    "admin@gms.local",
		"password": "admin123",
	}, "")
	if login.Code != http.StatusOK {
		t.Fatalf("login failed: %d body=%s", login.Code, login.Body.String())
	}
	token := decodeJSON[loginResp](t, login).Token

	createIndustry := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/industries", map[string]any{
		"name":   "Healthcare",
		"code":   "healthcare",
		"status": "draft",
	}, token)
	if createIndustry.Code != http.StatusCreated {
		t.Fatalf("create industry template failed: %d body=%s", createIndustry.Code, createIndustry.Body.String())
	}
	industryID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createIndustry).ID

	createProjectTemplate := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates", map[string]any{
		"industryTemplateId": industryID,
		"name":               "Clinic Setup",
		"code":               "clinic-setup",
		"status":             "draft",
	}, token)
	if createProjectTemplate.Code != http.StatusCreated {
		t.Fatalf("create project template failed: %d body=%s", createProjectTemplate.Code, createProjectTemplate.Body.String())
	}
	projectTemplateID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createProjectTemplate).ID

	createStageTemplate := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/stage-templates", map[string]any{
		"projectTemplateId": projectTemplateID,
		"name":              "Intake",
		"code":              "intake",
		"status":            "draft",
		"position":          1,
	}, token)
	if createStageTemplate.Code != http.StatusCreated {
		t.Fatalf("create stage template failed: %d body=%s", createStageTemplate.Code, createStageTemplate.Body.String())
	}
	stageTemplateID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createStageTemplate).ID

	createFormTemplate := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/form-templates", map[string]any{
		"stageTemplateId": stageTemplateID,
		"name":            "Patient Intake",
		"code":            "patient-intake",
		"status":          "draft",
		"position":        1,
	}, token)
	if createFormTemplate.Code != http.StatusCreated {
		t.Fatalf("create form template failed: %d body=%s", createFormTemplate.Code, createFormTemplate.Body.String())
	}
	formTemplateID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createFormTemplate).ID

	createFieldTemplate := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/field-templates", map[string]any{
		"formTemplateId": formTemplateID,
		"name":           "Patient Name",
		"code":           "patient-name",
		"status":         "draft",
		"position":       1,
		"widgetType":     "input",
	}, token)
	if createFieldTemplate.Code != http.StatusCreated {
		t.Fatalf("create field template failed: %d body=%s", createFieldTemplate.Code, createFieldTemplate.Body.String())
	}
	fieldTemplateID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createFieldTemplate).ID

	industries := doJSONRequest(t, h, http.MethodGet, "/api/tmpl/industries", nil, token)
	if industries.Code != http.StatusOK {
		t.Fatalf("list industries failed: %d body=%s", industries.Code, industries.Body.String())
	}
	industryItems := decodeJSON[listIndustryTemplatesResp](t, industries).Items
	if len(industryItems) == 0 || industryItems[0].ID != industryID {
		t.Fatalf("unexpected industries payload: %+v", industryItems)
	}

	projectTemplates := doJSONRequest(t, h, http.MethodGet, "/api/tmpl/project-templates", nil, token)
	if projectTemplates.Code != http.StatusOK {
		t.Fatalf("list project templates failed: %d body=%s", projectTemplates.Code, projectTemplates.Body.String())
	}
	projectItems := decodeJSON[listProjectTemplatesResp](t, projectTemplates).Items
	if len(projectItems) == 0 || projectItems[0].ID != projectTemplateID || projectItems[0].IndustryTemplateID != industryID {
		t.Fatalf("unexpected project template payload: %+v", projectItems)
	}

	stageTemplates := doJSONRequest(t, h, http.MethodGet, "/api/tmpl/stage-templates", nil, token)
	if stageTemplates.Code != http.StatusOK {
		t.Fatalf("list stage templates failed: %d body=%s", stageTemplates.Code, stageTemplates.Body.String())
	}
	stageItems := decodeJSON[listStageTemplatesResp](t, stageTemplates).Items
	if len(stageItems) == 0 || stageItems[0].ID != stageTemplateID || stageItems[0].ProjectTemplateID != projectTemplateID || stageItems[0].Position != 1 {
		t.Fatalf("unexpected stage template payload: %+v", stageItems)
	}

	formTemplates := doJSONRequest(t, h, http.MethodGet, "/api/tmpl/form-templates", nil, token)
	if formTemplates.Code != http.StatusOK {
		t.Fatalf("list form templates failed: %d body=%s", formTemplates.Code, formTemplates.Body.String())
	}
	formItems := decodeJSON[listFormTemplatesResp](t, formTemplates).Items
	if len(formItems) == 0 || formItems[0].ID != formTemplateID || formItems[0].StageTemplateID != stageTemplateID || formItems[0].Position != 1 {
		t.Fatalf("unexpected form template payload: %+v", formItems)
	}

	fieldTemplates := doJSONRequest(t, h, http.MethodGet, "/api/tmpl/field-templates", nil, token)
	if fieldTemplates.Code != http.StatusOK {
		t.Fatalf("list field templates failed: %d body=%s", fieldTemplates.Code, fieldTemplates.Body.String())
	}
	fieldItems := decodeJSON[listFieldTemplatesResp](t, fieldTemplates).Items
	if len(fieldItems) == 0 || fieldItems[0].ID != fieldTemplateID || fieldItems[0].FormTemplateID != formTemplateID || fieldItems[0].Position != 1 || fieldItems[0].WidgetType != "input" {
		t.Fatalf("unexpected field template payload: %+v", fieldItems)
	}
}
