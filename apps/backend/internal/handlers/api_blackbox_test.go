package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"gms/backend/internal/models"
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
		ID          int     `json:"id"`
		Name        string  `json:"name"`
		Status      string  `json:"status"`
		Version     int     `json:"version"`
		PublishedAt *string `json:"publishedAt"`
		PublishedBy *int    `json:"publishedBy"`
	} `json:"items"`
}

type listProjectTemplatesResp struct {
	Items []struct {
		ID                 int     `json:"id"`
		IndustryTemplateID int     `json:"industryTemplateId"`
		Status             string  `json:"status"`
		Version            int     `json:"version"`
		PublishedAt        *string `json:"publishedAt"`
		PublishedBy        *int    `json:"publishedBy"`
	} `json:"items"`
}

type listStageTemplatesResp struct {
	Items []struct {
		ID                int     `json:"id"`
		ProjectTemplateID int     `json:"projectTemplateId"`
		Position          int     `json:"position"`
		Status            string  `json:"status"`
		Version           int     `json:"version"`
		PublishedAt       *string `json:"publishedAt"`
		PublishedBy       *int    `json:"publishedBy"`
	} `json:"items"`
}

type listFormTemplatesResp struct {
	Items []struct {
		ID              int     `json:"id"`
		StageTemplateID int     `json:"stageTemplateId"`
		Position        int     `json:"position"`
		Status          string  `json:"status"`
		Version         int     `json:"version"`
		PublishedAt     *string `json:"publishedAt"`
		PublishedBy     *int    `json:"publishedBy"`
	} `json:"items"`
}

type listFieldTemplatesResp struct {
	Items []struct {
		ID             int     `json:"id"`
		FormTemplateID int     `json:"formTemplateId"`
		Position       int     `json:"position"`
		WidgetType     string  `json:"widgetType"`
		Status         string  `json:"status"`
		Version        int     `json:"version"`
		PublishedAt    *string `json:"publishedAt"`
		PublishedBy    *int    `json:"publishedBy"`
	} `json:"items"`
}

type templateStateResp struct {
	ID          int     `json:"id"`
	Status      string  `json:"status"`
	PublishedAt *string `json:"publishedAt"`
	PublishedBy *int    `json:"publishedBy"`
}

type instantiateProjectTemplateResp struct {
	Project struct {
		ID                     int    `json:"id"`
		Name                   string `json:"name"`
		Description            string `json:"description"`
		IndustryTemplateID     int    `json:"industryTemplateId"`
		ProjectTemplateID      int    `json:"projectTemplateId"`
		ProjectTemplateVersion int    `json:"projectTemplateVersion"`
		CreatedBy              int    `json:"createdBy"`
	} `json:"project"`
	Stages []struct {
		ID               int    `json:"id"`
		RuntimeProjectID int    `json:"runtimeProjectId"`
		StageTemplateID  int    `json:"stageTemplateId"`
		Name             string `json:"name"`
		Code             string `json:"code"`
		Description      string `json:"description"`
		Position         int    `json:"position"`
		Status           string `json:"status"`
	} `json:"stages"`
	Forms []struct {
		ID                    int    `json:"id"`
		RuntimeProjectID      int    `json:"runtimeProjectId"`
		RuntimeProjectStageID int    `json:"runtimeProjectStageId"`
		FormTemplateID        int    `json:"formTemplateId"`
		Name                  string `json:"name"`
		Code                  string `json:"code"`
		Description           string `json:"description"`
		Position              int    `json:"position"`
	} `json:"forms"`
	Fields []struct {
		ID                   int     `json:"id"`
		RuntimeProjectID     int     `json:"runtimeProjectId"`
		RuntimeProjectFormID int     `json:"runtimeProjectFormId"`
		FormFieldTemplateID  int     `json:"formFieldTemplateId"`
		Name                 string  `json:"name"`
		Code                 string  `json:"code"`
		Description          string  `json:"description"`
		Position             int     `json:"position"`
		WidgetType           string  `json:"widgetType"`
		ValueText            *string `json:"valueText"`
	} `json:"fields"`
}

type nextProjectTemplateVersionResp struct {
	ProjectTemplate struct {
		ID                 int     `json:"id"`
		IndustryTemplateID int     `json:"industryTemplateId"`
		Name               string  `json:"name"`
		Code               string  `json:"code"`
		Description        string  `json:"description"`
		Version            int     `json:"version"`
		Status             string  `json:"status"`
		PublishedAt        *string `json:"publishedAt"`
		PublishedBy        *int    `json:"publishedBy"`
	} `json:"projectTemplate"`
	Stages []struct {
		ID                int     `json:"id"`
		ProjectTemplateID int     `json:"projectTemplateId"`
		Name              string  `json:"name"`
		Code              string  `json:"code"`
		Description       string  `json:"description"`
		Version           int     `json:"version"`
		Status            string  `json:"status"`
		Position          int     `json:"position"`
		PublishedAt       *string `json:"publishedAt"`
		PublishedBy       *int    `json:"publishedBy"`
	} `json:"stages"`
	Forms []struct {
		ID              int     `json:"id"`
		StageTemplateID int     `json:"stageTemplateId"`
		Name            string  `json:"name"`
		Code            string  `json:"code"`
		Description     string  `json:"description"`
		Version         int     `json:"version"`
		Status          string  `json:"status"`
		Position        int     `json:"position"`
		PublishedAt     *string `json:"publishedAt"`
		PublishedBy     *int    `json:"publishedBy"`
	} `json:"forms"`
	Fields []struct {
		ID             int     `json:"id"`
		FormTemplateID int     `json:"formTemplateId"`
		Name           string  `json:"name"`
		Code           string  `json:"code"`
		Description    string  `json:"description"`
		Version        int     `json:"version"`
		Status         string  `json:"status"`
		Position       int     `json:"position"`
		WidgetType     string  `json:"widgetType"`
		PublishedAt    *string `json:"publishedAt"`
		PublishedBy    *int    `json:"publishedBy"`
	} `json:"fields"`
}

type projectTemplateLifecycleResp struct {
	ProjectTemplate struct {
		ID     int    `json:"id"`
		Status string `json:"status"`
	} `json:"projectTemplate"`
	IndustryTemplate struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Code    string `json:"code"`
		Status  string `json:"status"`
		Version int    `json:"version"`
	} `json:"industryTemplate"`
	Counts struct {
		PublishedStages int `json:"publishedStages"`
		PublishedForms  int `json:"publishedForms"`
		PublishedFields int `json:"publishedFields"`
		RuntimeProjects int `json:"runtimeProjects"`
	} `json:"counts"`
	RuntimeByStageStatus struct {
		Pending int `json:"pending"`
		Active  int `json:"active"`
		Done    int `json:"done"`
	} `json:"runtimeByStageStatus"`
	Guidance string `json:"guidance"`
}

type errorResp struct {
	Error string `json:"error"`
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

func TestBlackbox_TemplateEngineProjectTemplateListFilters(t *testing.T) {
	h, _, _ := setupTestRouter(t)

	login := doJSONRequest(t, h, http.MethodPost, "/api/auth/login", map[string]any{
		"email":    "admin@gms.local",
		"password": "admin123",
	}, "")
	if login.Code != http.StatusOK {
		t.Fatalf("login failed: %d body=%s", login.Code, login.Body.String())
	}
	token := decodeJSON[loginResp](t, login).Token

	createIndustryOne := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/industries", map[string]any{
		"name": "Industry One",
		"code": "industry-one",
	}, token)
	if createIndustryOne.Code != http.StatusCreated {
		t.Fatalf("create industry one failed: %d body=%s", createIndustryOne.Code, createIndustryOne.Body.String())
	}
	industryOneID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createIndustryOne).ID

	createIndustryTwo := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/industries", map[string]any{
		"name": "Industry Two",
		"code": "industry-two",
	}, token)
	if createIndustryTwo.Code != http.StatusCreated {
		t.Fatalf("create industry two failed: %d body=%s", createIndustryTwo.Code, createIndustryTwo.Body.String())
	}
	industryTwoID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createIndustryTwo).ID

	createProjectOne := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates", map[string]any{
		"industryTemplateId": industryOneID,
		"name":               "Draft V1",
		"code":               "draft-v1",
		"status":             "draft",
		"version":            1,
	}, token)
	if createProjectOne.Code != http.StatusCreated {
		t.Fatalf("create project one failed: %d body=%s", createProjectOne.Code, createProjectOne.Body.String())
	}
	projectOneID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createProjectOne).ID

	createProjectTwo := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates", map[string]any{
		"industryTemplateId": industryOneID,
		"name":               "Published V2",
		"code":               "published-v2",
		"status":             "published",
		"version":            2,
	}, token)
	if createProjectTwo.Code != http.StatusCreated {
		t.Fatalf("create project two failed: %d body=%s", createProjectTwo.Code, createProjectTwo.Body.String())
	}

	createProjectThree := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates", map[string]any{
		"industryTemplateId": industryTwoID,
		"name":               "Other Industry",
		"code":               "other-industry",
		"status":             "draft",
		"version":            1,
	}, token)
	if createProjectThree.Code != http.StatusCreated {
		t.Fatalf("create project three failed: %d body=%s", createProjectThree.Code, createProjectThree.Body.String())
	}

	filtered := doJSONRequest(t, h, http.MethodGet, "/api/tmpl/project-templates?industryTemplateId="+itoa(industryOneID)+"&status=draft&version=1", nil, token)
	if filtered.Code != http.StatusOK {
		t.Fatalf("list project templates with filters failed: %d body=%s", filtered.Code, filtered.Body.String())
	}
	filteredItems := decodeJSON[listProjectTemplatesResp](t, filtered).Items
	if len(filteredItems) != 1 || filteredItems[0].ID != projectOneID {
		t.Fatalf("unexpected filtered project templates payload: %+v", filteredItems)
	}

	invalidStatus := doJSONRequest(t, h, http.MethodGet, "/api/tmpl/project-templates?status=invalid", nil, token)
	if invalidStatus.Code != http.StatusBadRequest {
		t.Fatalf("expected invalid status 400, got %d body=%s", invalidStatus.Code, invalidStatus.Body.String())
	}

	invalidVersion := doJSONRequest(t, h, http.MethodGet, "/api/tmpl/project-templates?version=0", nil, token)
	if invalidVersion.Code != http.StatusBadRequest {
		t.Fatalf("expected invalid version 400, got %d body=%s", invalidVersion.Code, invalidVersion.Body.String())
	}

	invalidParent := doJSONRequest(t, h, http.MethodGet, "/api/tmpl/project-templates?industryTemplateId=0", nil, token)
	if invalidParent.Code != http.StatusBadRequest {
		t.Fatalf("expected invalid industryTemplateId 400, got %d body=%s", invalidParent.Code, invalidParent.Body.String())
	}
}

func TestBlackbox_TemplateEnginePublishChainSuccess(t *testing.T) {
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
	userID := loginPayload.User.ID

	createIndustry := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/industries", map[string]any{
		"name": "Manufacturing",
		"code": "manufacturing",
	}, token)
	if createIndustry.Code != http.StatusCreated {
		t.Fatalf("create industry failed: %d body=%s", createIndustry.Code, createIndustry.Body.String())
	}
	industryID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createIndustry).ID

	createProject := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates", map[string]any{
		"industryTemplateId": industryID,
		"name":               "Factory Rollout",
		"code":               "factory-rollout",
	}, token)
	if createProject.Code != http.StatusCreated {
		t.Fatalf("create project template failed: %d body=%s", createProject.Code, createProject.Body.String())
	}
	projectID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createProject).ID

	createStage := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/stage-templates", map[string]any{
		"projectTemplateId": projectID,
		"name":              "Planning",
		"code":              "planning",
		"position":          1,
	}, token)
	if createStage.Code != http.StatusCreated {
		t.Fatalf("create stage template failed: %d body=%s", createStage.Code, createStage.Body.String())
	}
	stageID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createStage).ID

	createForm := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/form-templates", map[string]any{
		"stageTemplateId": stageID,
		"name":            "Checklist",
		"code":            "checklist",
		"position":        1,
	}, token)
	if createForm.Code != http.StatusCreated {
		t.Fatalf("create form template failed: %d body=%s", createForm.Code, createForm.Body.String())
	}
	formID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createForm).ID

	createField := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/field-templates", map[string]any{
		"formTemplateId": formID,
		"name":           "Approver",
		"code":           "approver",
		"position":       1,
		"widgetType":     "input",
	}, token)
	if createField.Code != http.StatusCreated {
		t.Fatalf("create field template failed: %d body=%s", createField.Code, createField.Body.String())
	}
	fieldID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createField).ID

	assertPublished := func(resp templateStateResp, expectedID int) {
		t.Helper()
		if resp.ID != expectedID {
			t.Fatalf("unexpected template id: want=%d got=%d", expectedID, resp.ID)
		}
		if resp.Status != "published" {
			t.Fatalf("expected status published, got %s", resp.Status)
		}
		if resp.PublishedAt == nil || *resp.PublishedAt == "" {
			t.Fatalf("expected publishedAt to be set")
		}
		if resp.PublishedBy == nil || *resp.PublishedBy != userID {
			t.Fatalf("expected publishedBy=%d, got %+v", userID, resp.PublishedBy)
		}
	}

	publishIndustry := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/industries/"+itoa(industryID)+"/publish", nil, token)
	if publishIndustry.Code != http.StatusOK {
		t.Fatalf("publish industry failed: %d body=%s", publishIndustry.Code, publishIndustry.Body.String())
	}
	assertPublished(decodeJSON[templateStateResp](t, publishIndustry), industryID)

	publishProject := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates/"+itoa(projectID)+"/publish", nil, token)
	if publishProject.Code != http.StatusOK {
		t.Fatalf("publish project failed: %d body=%s", publishProject.Code, publishProject.Body.String())
	}
	assertPublished(decodeJSON[templateStateResp](t, publishProject), projectID)

	publishStage := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/stage-templates/"+itoa(stageID)+"/publish", nil, token)
	if publishStage.Code != http.StatusOK {
		t.Fatalf("publish stage failed: %d body=%s", publishStage.Code, publishStage.Body.String())
	}
	assertPublished(decodeJSON[templateStateResp](t, publishStage), stageID)

	publishForm := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/form-templates/"+itoa(formID)+"/publish", nil, token)
	if publishForm.Code != http.StatusOK {
		t.Fatalf("publish form failed: %d body=%s", publishForm.Code, publishForm.Body.String())
	}
	assertPublished(decodeJSON[templateStateResp](t, publishForm), formID)

	publishField := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/field-templates/"+itoa(fieldID)+"/publish", nil, token)
	if publishField.Code != http.StatusOK {
		t.Fatalf("publish field failed: %d body=%s", publishField.Code, publishField.Body.String())
	}
	assertPublished(decodeJSON[templateStateResp](t, publishField), fieldID)
}

func TestBlackbox_TemplateEnginePublishBlockedWhenParentDraft(t *testing.T) {
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
		"name": "Energy",
		"code": "energy",
	}, token)
	if createIndustry.Code != http.StatusCreated {
		t.Fatalf("create industry failed: %d body=%s", createIndustry.Code, createIndustry.Body.String())
	}
	industryID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createIndustry).ID

	createProject := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates", map[string]any{
		"industryTemplateId": industryID,
		"name":               "Solar Setup",
		"code":               "solar-setup",
	}, token)
	if createProject.Code != http.StatusCreated {
		t.Fatalf("create project failed: %d body=%s", createProject.Code, createProject.Body.String())
	}
	projectID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createProject).ID

	publishProject := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates/"+itoa(projectID)+"/publish", nil, token)
	if publishProject.Code != http.StatusBadRequest {
		t.Fatalf("expected publish project 400 when parent is draft, got %d body=%s", publishProject.Code, publishProject.Body.String())
	}
	errPayload := decodeJSON[errorResp](t, publishProject)
	if !strings.Contains(errPayload.Error, "industry template must be published first") {
		t.Fatalf("unexpected hierarchy error message: %s", errPayload.Error)
	}
}

func TestBlackbox_TemplateEngineUnpublishBlockedWhenPublishedChildExists(t *testing.T) {
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
		"name": "Logistics",
		"code": "logistics",
	}, token)
	if createIndustry.Code != http.StatusCreated {
		t.Fatalf("create industry failed: %d body=%s", createIndustry.Code, createIndustry.Body.String())
	}
	industryID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createIndustry).ID

	createProject := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates", map[string]any{
		"industryTemplateId": industryID,
		"name":               "Warehouse Build",
		"code":               "warehouse-build",
	}, token)
	if createProject.Code != http.StatusCreated {
		t.Fatalf("create project failed: %d body=%s", createProject.Code, createProject.Body.String())
	}
	projectID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createProject).ID

	publishIndustry := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/industries/"+itoa(industryID)+"/publish", nil, token)
	if publishIndustry.Code != http.StatusOK {
		t.Fatalf("publish industry failed: %d body=%s", publishIndustry.Code, publishIndustry.Body.String())
	}

	publishProject := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates/"+itoa(projectID)+"/publish", nil, token)
	if publishProject.Code != http.StatusOK {
		t.Fatalf("publish project failed: %d body=%s", publishProject.Code, publishProject.Body.String())
	}

	unpublishIndustry := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/industries/"+itoa(industryID)+"/unpublish", nil, token)
	if unpublishIndustry.Code != http.StatusBadRequest {
		t.Fatalf("expected unpublish industry 400 while published child exists, got %d body=%s", unpublishIndustry.Code, unpublishIndustry.Body.String())
	}
	errPayload := decodeJSON[errorResp](t, unpublishIndustry)
	if !strings.Contains(errPayload.Error, "published project templates exist") {
		t.Fatalf("unexpected hierarchy error message: %s", errPayload.Error)
	}
}

func TestBlackbox_TemplateEngineInstantiateSuccess(t *testing.T) {
	h, db, _ := setupTestRouter(t)

	login := doJSONRequest(t, h, http.MethodPost, "/api/auth/login", map[string]any{
		"email":    "admin@gms.local",
		"password": "admin123",
	}, "")
	if login.Code != http.StatusOK {
		t.Fatalf("login failed: %d body=%s", login.Code, login.Body.String())
	}
	loginPayload := decodeJSON[loginResp](t, login)
	token := loginPayload.Token
	userID := loginPayload.User.ID

	createIndustry := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/industries", map[string]any{
		"name": "Construction",
		"code": "construction",
	}, token)
	if createIndustry.Code != http.StatusCreated {
		t.Fatalf("create industry failed: %d body=%s", createIndustry.Code, createIndustry.Body.String())
	}
	industryID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createIndustry).ID

	createProject := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates", map[string]any{
		"industryTemplateId": industryID,
		"name":               "Site Rollout",
		"code":               "site-rollout",
	}, token)
	if createProject.Code != http.StatusCreated {
		t.Fatalf("create project failed: %d body=%s", createProject.Code, createProject.Body.String())
	}
	projectID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createProject).ID

	createStageOne := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/stage-templates", map[string]any{
		"projectTemplateId": projectID,
		"name":              "Planning",
		"code":              "planning",
		"position":          1,
	}, token)
	if createStageOne.Code != http.StatusCreated {
		t.Fatalf("create stage one failed: %d body=%s", createStageOne.Code, createStageOne.Body.String())
	}
	stageOneID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createStageOne).ID

	createStageTwo := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/stage-templates", map[string]any{
		"projectTemplateId": projectID,
		"name":              "Execution",
		"code":              "execution",
		"position":          2,
	}, token)
	if createStageTwo.Code != http.StatusCreated {
		t.Fatalf("create stage two failed: %d body=%s", createStageTwo.Code, createStageTwo.Body.String())
	}
	stageTwoID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createStageTwo).ID

	createFormOne := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/form-templates", map[string]any{
		"stageTemplateId": stageOneID,
		"name":            "Plan Form",
		"code":            "plan-form",
		"position":        1,
	}, token)
	if createFormOne.Code != http.StatusCreated {
		t.Fatalf("create form one failed: %d body=%s", createFormOne.Code, createFormOne.Body.String())
	}
	formOneID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createFormOne).ID

	createFormTwo := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/form-templates", map[string]any{
		"stageTemplateId": stageTwoID,
		"name":            "Execution Form",
		"code":            "execution-form",
		"position":        1,
	}, token)
	if createFormTwo.Code != http.StatusCreated {
		t.Fatalf("create form two failed: %d body=%s", createFormTwo.Code, createFormTwo.Body.String())
	}
	formTwoID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createFormTwo).ID

	createFieldOne := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/field-templates", map[string]any{
		"formTemplateId": formOneID,
		"name":           "Site Name",
		"code":           "site-name",
		"position":       1,
		"widgetType":     "input",
	}, token)
	if createFieldOne.Code != http.StatusCreated {
		t.Fatalf("create field one failed: %d body=%s", createFieldOne.Code, createFieldOne.Body.String())
	}
	fieldOneID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createFieldOne).ID

	createFieldTwo := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/field-templates", map[string]any{
		"formTemplateId": formTwoID,
		"name":           "Start Date",
		"code":           "start-date",
		"position":       1,
		"widgetType":     "date",
	}, token)
	if createFieldTwo.Code != http.StatusCreated {
		t.Fatalf("create field two failed: %d body=%s", createFieldTwo.Code, createFieldTwo.Body.String())
	}
	fieldTwoID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createFieldTwo).ID

	publishIndustry := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/industries/"+itoa(industryID)+"/publish", nil, token)
	if publishIndustry.Code != http.StatusOK {
		t.Fatalf("publish industry failed: %d body=%s", publishIndustry.Code, publishIndustry.Body.String())
	}

	publishProject := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates/"+itoa(projectID)+"/publish", nil, token)
	if publishProject.Code != http.StatusOK {
		t.Fatalf("publish project failed: %d body=%s", publishProject.Code, publishProject.Body.String())
	}

	publishStageOne := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/stage-templates/"+itoa(stageOneID)+"/publish", nil, token)
	if publishStageOne.Code != http.StatusOK {
		t.Fatalf("publish stage one failed: %d body=%s", publishStageOne.Code, publishStageOne.Body.String())
	}

	publishStageTwo := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/stage-templates/"+itoa(stageTwoID)+"/publish", nil, token)
	if publishStageTwo.Code != http.StatusOK {
		t.Fatalf("publish stage two failed: %d body=%s", publishStageTwo.Code, publishStageTwo.Body.String())
	}

	publishFormOne := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/form-templates/"+itoa(formOneID)+"/publish", nil, token)
	if publishFormOne.Code != http.StatusOK {
		t.Fatalf("publish form one failed: %d body=%s", publishFormOne.Code, publishFormOne.Body.String())
	}

	publishFormTwo := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/form-templates/"+itoa(formTwoID)+"/publish", nil, token)
	if publishFormTwo.Code != http.StatusOK {
		t.Fatalf("publish form two failed: %d body=%s", publishFormTwo.Code, publishFormTwo.Body.String())
	}

	publishFieldOne := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/field-templates/"+itoa(fieldOneID)+"/publish", nil, token)
	if publishFieldOne.Code != http.StatusOK {
		t.Fatalf("publish field one failed: %d body=%s", publishFieldOne.Code, publishFieldOne.Body.String())
	}

	publishFieldTwo := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/field-templates/"+itoa(fieldTwoID)+"/publish", nil, token)
	if publishFieldTwo.Code != http.StatusOK {
		t.Fatalf("publish field two failed: %d body=%s", publishFieldTwo.Code, publishFieldTwo.Body.String())
	}

	instantiate := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates/"+itoa(projectID)+"/instantiate", map[string]any{
		"name":        "Runtime Site Rollout",
		"description": "runtime project generated in blackbox test",
	}, token)
	if instantiate.Code != http.StatusCreated {
		t.Fatalf("instantiate failed: %d body=%s", instantiate.Code, instantiate.Body.String())
	}
	payload := decodeJSON[instantiateProjectTemplateResp](t, instantiate)

	if payload.Project.ID <= 0 {
		t.Fatalf("expected runtime project id to be set")
	}
	if payload.Project.Name != "Runtime Site Rollout" {
		t.Fatalf("unexpected runtime project name: %s", payload.Project.Name)
	}
	if payload.Project.Description != "runtime project generated in blackbox test" {
		t.Fatalf("unexpected runtime project description: %s", payload.Project.Description)
	}
	if payload.Project.ProjectTemplateID != projectID {
		t.Fatalf("expected projectTemplateId=%d, got %d", projectID, payload.Project.ProjectTemplateID)
	}
	if payload.Project.IndustryTemplateID != industryID {
		t.Fatalf("expected industryTemplateId=%d, got %d", industryID, payload.Project.IndustryTemplateID)
	}
	if payload.Project.ProjectTemplateVersion != 1 {
		t.Fatalf("expected projectTemplateVersion=1, got %d", payload.Project.ProjectTemplateVersion)
	}
	if payload.Project.CreatedBy != userID {
		t.Fatalf("expected createdBy=%d, got %d", userID, payload.Project.CreatedBy)
	}

	if len(payload.Stages) == 0 {
		t.Fatalf("expected runtime stages to be created")
	}
	if len(payload.Forms) == 0 {
		t.Fatalf("expected runtime forms to be created")
	}
	if len(payload.Fields) == 0 {
		t.Fatalf("expected runtime fields to be created")
	}
	if payload.Stages[0].Status != "active" {
		t.Fatalf("expected first stage active, got %s", payload.Stages[0].Status)
	}
	for i := 1; i < len(payload.Stages); i++ {
		if payload.Stages[i].Status != "pending" {
			t.Fatalf("expected stage %d pending, got %s", i, payload.Stages[i].Status)
		}
	}

	var runtimeProject models.RuntimeProject
	if err := db.First(&runtimeProject, payload.Project.ID).Error; err != nil {
		t.Fatalf("failed to load runtime project from db: %v", err)
	}

	var stageCount int64
	if err := db.Model(&models.RuntimeProjectStage{}).Where("runtime_project_id = ?", payload.Project.ID).Count(&stageCount).Error; err != nil {
		t.Fatalf("failed counting runtime stages: %v", err)
	}
	if int(stageCount) != len(payload.Stages) {
		t.Fatalf("expected stage count=%d, got %d", len(payload.Stages), stageCount)
	}

	var formCount int64
	if err := db.Model(&models.RuntimeProjectForm{}).Where("runtime_project_id = ?", payload.Project.ID).Count(&formCount).Error; err != nil {
		t.Fatalf("failed counting runtime forms: %v", err)
	}
	if int(formCount) != len(payload.Forms) {
		t.Fatalf("expected form count=%d, got %d", len(payload.Forms), formCount)
	}

	var fieldCount int64
	if err := db.Model(&models.RuntimeProjectField{}).Where("runtime_project_id = ?", payload.Project.ID).Count(&fieldCount).Error; err != nil {
		t.Fatalf("failed counting runtime fields: %v", err)
	}
	if int(fieldCount) != len(payload.Fields) {
		t.Fatalf("expected field count=%d, got %d", len(payload.Fields), fieldCount)
	}
}

func TestBlackbox_TemplateEngineInstantiateFailsWhenProjectTemplateDraft(t *testing.T) {
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
		"name": "Retail",
		"code": "retail",
	}, token)
	if createIndustry.Code != http.StatusCreated {
		t.Fatalf("create industry failed: %d body=%s", createIndustry.Code, createIndustry.Body.String())
	}
	industryID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createIndustry).ID

	publishIndustry := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/industries/"+itoa(industryID)+"/publish", nil, token)
	if publishIndustry.Code != http.StatusOK {
		t.Fatalf("publish industry failed: %d body=%s", publishIndustry.Code, publishIndustry.Body.String())
	}

	createProject := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates", map[string]any{
		"industryTemplateId": industryID,
		"name":               "Store Rollout",
		"code":               "store-rollout",
		"status":             "draft",
	}, token)
	if createProject.Code != http.StatusCreated {
		t.Fatalf("create project failed: %d body=%s", createProject.Code, createProject.Body.String())
	}
	projectID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createProject).ID

	instantiate := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates/"+itoa(projectID)+"/instantiate", map[string]any{
		"name": "Runtime Store Rollout",
	}, token)
	if instantiate.Code != http.StatusBadRequest {
		t.Fatalf("expected instantiate draft project 400, got %d body=%s", instantiate.Code, instantiate.Body.String())
	}

	errPayload := decodeJSON[errorResp](t, instantiate)
	if !strings.Contains(errPayload.Error, "template must be published") {
		t.Fatalf("unexpected draft error message: %s", errPayload.Error)
	}
}

func TestBlackbox_TemplateEngineInstantiateFailsWhenNoPublishedStages(t *testing.T) {
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
		"name": "Hospitality",
		"code": "hospitality",
	}, token)
	if createIndustry.Code != http.StatusCreated {
		t.Fatalf("create industry failed: %d body=%s", createIndustry.Code, createIndustry.Body.String())
	}
	industryID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createIndustry).ID

	createProject := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates", map[string]any{
		"industryTemplateId": industryID,
		"name":               "Hotel Launch",
		"code":               "hotel-launch",
	}, token)
	if createProject.Code != http.StatusCreated {
		t.Fatalf("create project failed: %d body=%s", createProject.Code, createProject.Body.String())
	}
	projectID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createProject).ID

	createStage := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/stage-templates", map[string]any{
		"projectTemplateId": projectID,
		"name":              "Draft Stage",
		"code":              "draft-stage",
		"status":            "draft",
		"position":          1,
	}, token)
	if createStage.Code != http.StatusCreated {
		t.Fatalf("create stage failed: %d body=%s", createStage.Code, createStage.Body.String())
	}

	publishIndustry := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/industries/"+itoa(industryID)+"/publish", nil, token)
	if publishIndustry.Code != http.StatusOK {
		t.Fatalf("publish industry failed: %d body=%s", publishIndustry.Code, publishIndustry.Body.String())
	}

	publishProject := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates/"+itoa(projectID)+"/publish", nil, token)
	if publishProject.Code != http.StatusOK {
		t.Fatalf("publish project failed: %d body=%s", publishProject.Code, publishProject.Body.String())
	}

	instantiate := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates/"+itoa(projectID)+"/instantiate", map[string]any{
		"name": "Runtime Hotel Launch",
	}, token)
	if instantiate.Code != http.StatusBadRequest {
		t.Fatalf("expected instantiate no-published-stage 400, got %d body=%s", instantiate.Code, instantiate.Body.String())
	}

	errPayload := decodeJSON[errorResp](t, instantiate)
	if !strings.Contains(errPayload.Error, "no published stage templates found") {
		t.Fatalf("unexpected stage error message: %s", errPayload.Error)
	}
}

func TestBlackbox_TemplateEngineNextVersionClonesHierarchyAsDraft(t *testing.T) {
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
		"name": "Governance Industry",
		"code": "governance-industry",
	}, token)
	if createIndustry.Code != http.StatusCreated {
		t.Fatalf("create industry failed: %d body=%s", createIndustry.Code, createIndustry.Body.String())
	}
	industryID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createIndustry).ID

	publishIndustry := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/industries/"+itoa(industryID)+"/publish", nil, token)
	if publishIndustry.Code != http.StatusOK {
		t.Fatalf("publish industry failed: %d body=%s", publishIndustry.Code, publishIndustry.Body.String())
	}

	createProject := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates", map[string]any{
		"industryTemplateId": industryID,
		"name":               "Ops Template",
		"code":               "ops-template",
		"description":        "baseline template for ops",
	}, token)
	if createProject.Code != http.StatusCreated {
		t.Fatalf("create project template failed: %d body=%s", createProject.Code, createProject.Body.String())
	}
	projectID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createProject).ID

	createStageOne := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/stage-templates", map[string]any{
		"projectTemplateId": projectID,
		"name":              "Design",
		"code":              "design",
		"position":          1,
	}, token)
	if createStageOne.Code != http.StatusCreated {
		t.Fatalf("create stage one failed: %d body=%s", createStageOne.Code, createStageOne.Body.String())
	}
	stageOneID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createStageOne).ID

	createStageTwo := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/stage-templates", map[string]any{
		"projectTemplateId": projectID,
		"name":              "Delivery",
		"code":              "delivery",
		"position":          2,
	}, token)
	if createStageTwo.Code != http.StatusCreated {
		t.Fatalf("create stage two failed: %d body=%s", createStageTwo.Code, createStageTwo.Body.String())
	}
	stageTwoID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createStageTwo).ID

	createFormOne := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/form-templates", map[string]any{
		"stageTemplateId": stageOneID,
		"name":            "Design Form",
		"code":            "design-form",
		"position":        1,
	}, token)
	if createFormOne.Code != http.StatusCreated {
		t.Fatalf("create form one failed: %d body=%s", createFormOne.Code, createFormOne.Body.String())
	}
	formOneID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createFormOne).ID

	createFormTwo := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/form-templates", map[string]any{
		"stageTemplateId": stageTwoID,
		"name":            "Delivery Form",
		"code":            "delivery-form",
		"position":        2,
	}, token)
	if createFormTwo.Code != http.StatusCreated {
		t.Fatalf("create form two failed: %d body=%s", createFormTwo.Code, createFormTwo.Body.String())
	}
	formTwoID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createFormTwo).ID

	createFieldOne := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/field-templates", map[string]any{
		"formTemplateId": formOneID,
		"name":           "Owner",
		"code":           "owner",
		"position":       1,
		"widgetType":     "input",
	}, token)
	if createFieldOne.Code != http.StatusCreated {
		t.Fatalf("create field one failed: %d body=%s", createFieldOne.Code, createFieldOne.Body.String())
	}
	fieldOneID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createFieldOne).ID

	createFieldTwo := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/field-templates", map[string]any{
		"formTemplateId": formTwoID,
		"name":           "Go Live Date",
		"code":           "go-live-date",
		"position":       2,
		"widgetType":     "date",
	}, token)
	if createFieldTwo.Code != http.StatusCreated {
		t.Fatalf("create field two failed: %d body=%s", createFieldTwo.Code, createFieldTwo.Body.String())
	}

	publishProject := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates/"+itoa(projectID)+"/publish", nil, token)
	if publishProject.Code != http.StatusOK {
		t.Fatalf("publish project failed: %d body=%s", publishProject.Code, publishProject.Body.String())
	}

	publishStageOne := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/stage-templates/"+itoa(stageOneID)+"/publish", nil, token)
	if publishStageOne.Code != http.StatusOK {
		t.Fatalf("publish stage one failed: %d body=%s", publishStageOne.Code, publishStageOne.Body.String())
	}

	publishFormOne := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/form-templates/"+itoa(formOneID)+"/publish", nil, token)
	if publishFormOne.Code != http.StatusOK {
		t.Fatalf("publish form one failed: %d body=%s", publishFormOne.Code, publishFormOne.Body.String())
	}

	publishFieldOne := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/field-templates/"+itoa(fieldOneID)+"/publish", nil, token)
	if publishFieldOne.Code != http.StatusOK {
		t.Fatalf("publish field one failed: %d body=%s", publishFieldOne.Code, publishFieldOne.Body.String())
	}

	nextVersion := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates/"+itoa(projectID)+"/next-version", nil, token)
	if nextVersion.Code != http.StatusCreated {
		t.Fatalf("next-version failed: %d body=%s", nextVersion.Code, nextVersion.Body.String())
	}
	payload := decodeJSON[nextProjectTemplateVersionResp](t, nextVersion)

	if payload.ProjectTemplate.ID <= 0 {
		t.Fatalf("expected cloned project template id to be set")
	}
	if payload.ProjectTemplate.ID == projectID {
		t.Fatalf("expected cloned project template id to differ from source")
	}
	if payload.ProjectTemplate.IndustryTemplateID != industryID {
		t.Fatalf("expected cloned industryTemplateId=%d, got %d", industryID, payload.ProjectTemplate.IndustryTemplateID)
	}
	if payload.ProjectTemplate.Code != "ops-template" {
		t.Fatalf("expected cloned code to match source, got %s", payload.ProjectTemplate.Code)
	}
	if payload.ProjectTemplate.Description != "baseline template for ops" {
		t.Fatalf("expected cloned description to match source, got %s", payload.ProjectTemplate.Description)
	}
	if payload.ProjectTemplate.Version != 2 {
		t.Fatalf("expected cloned version=2, got %d", payload.ProjectTemplate.Version)
	}
	if payload.ProjectTemplate.Status != "draft" {
		t.Fatalf("expected cloned project status=draft, got %s", payload.ProjectTemplate.Status)
	}
	if payload.ProjectTemplate.PublishedAt != nil || payload.ProjectTemplate.PublishedBy != nil {
		t.Fatalf("expected cloned project to clear publish metadata")
	}

	if len(payload.Stages) != 2 {
		t.Fatalf("expected 2 cloned stages, got %d", len(payload.Stages))
	}
	if len(payload.Forms) != 2 {
		t.Fatalf("expected 2 cloned forms, got %d", len(payload.Forms))
	}
	if len(payload.Fields) != 2 {
		t.Fatalf("expected 2 cloned fields, got %d", len(payload.Fields))
	}

	clonedStageIDs := make(map[int]struct{}, len(payload.Stages))
	for _, stage := range payload.Stages {
		if stage.ProjectTemplateID != payload.ProjectTemplate.ID {
			t.Fatalf("expected cloned stage projectTemplateId=%d, got %d", payload.ProjectTemplate.ID, stage.ProjectTemplateID)
		}
		if stage.Status != "draft" {
			t.Fatalf("expected cloned stage status=draft, got %s", stage.Status)
		}
		if stage.PublishedAt != nil || stage.PublishedBy != nil {
			t.Fatalf("expected cloned stage to clear publish metadata")
		}
		clonedStageIDs[stage.ID] = struct{}{}
	}

	clonedFormIDs := make(map[int]struct{}, len(payload.Forms))
	for _, form := range payload.Forms {
		if _, exists := clonedStageIDs[form.StageTemplateID]; !exists {
			t.Fatalf("expected cloned form to point to cloned stage, got stageTemplateId=%d", form.StageTemplateID)
		}
		if form.Status != "draft" {
			t.Fatalf("expected cloned form status=draft, got %s", form.Status)
		}
		if form.PublishedAt != nil || form.PublishedBy != nil {
			t.Fatalf("expected cloned form to clear publish metadata")
		}
		clonedFormIDs[form.ID] = struct{}{}
	}

	for _, field := range payload.Fields {
		if _, exists := clonedFormIDs[field.FormTemplateID]; !exists {
			t.Fatalf("expected cloned field to point to cloned form, got formTemplateId=%d", field.FormTemplateID)
		}
		if field.Status != "draft" {
			t.Fatalf("expected cloned field status=draft, got %s", field.Status)
		}
		if field.PublishedAt != nil || field.PublishedBy != nil {
			t.Fatalf("expected cloned field to clear publish metadata")
		}
	}
}

func TestBlackbox_TemplateEngineNextVersionFailsWhenSourceProjectTemplateDraft(t *testing.T) {
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
		"name": "Governance Draft Industry",
		"code": "governance-draft-industry",
	}, token)
	if createIndustry.Code != http.StatusCreated {
		t.Fatalf("create industry failed: %d body=%s", createIndustry.Code, createIndustry.Body.String())
	}
	industryID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createIndustry).ID

	createProject := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates", map[string]any{
		"industryTemplateId": industryID,
		"name":               "Draft Template",
		"code":               "draft-template",
		"status":             "draft",
	}, token)
	if createProject.Code != http.StatusCreated {
		t.Fatalf("create project failed: %d body=%s", createProject.Code, createProject.Body.String())
	}
	projectID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createProject).ID

	nextVersion := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates/"+itoa(projectID)+"/next-version", nil, token)
	if nextVersion.Code != http.StatusBadRequest {
		t.Fatalf("expected next-version on draft source to return 400, got %d body=%s", nextVersion.Code, nextVersion.Body.String())
	}

	errPayload := decodeJSON[errorResp](t, nextVersion)
	if !strings.Contains(errPayload.Error, "must be published") {
		t.Fatalf("unexpected draft source error message: %s", errPayload.Error)
	}
}

func TestBlackbox_TemplateEngineProjectTemplateLifecycleSummary(t *testing.T) {
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
		"name": "Lifecycle Industry",
		"code": "lifecycle-industry",
	}, token)
	if createIndustry.Code != http.StatusCreated {
		t.Fatalf("create industry failed: %d body=%s", createIndustry.Code, createIndustry.Body.String())
	}
	industryID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createIndustry).ID

	publishIndustry := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/industries/"+itoa(industryID)+"/publish", nil, token)
	if publishIndustry.Code != http.StatusOK {
		t.Fatalf("publish industry failed: %d body=%s", publishIndustry.Code, publishIndustry.Body.String())
	}

	createProject := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates", map[string]any{
		"industryTemplateId": industryID,
		"name":               "Lifecycle Project",
		"code":               "lifecycle-project",
	}, token)
	if createProject.Code != http.StatusCreated {
		t.Fatalf("create project failed: %d body=%s", createProject.Code, createProject.Body.String())
	}
	projectID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createProject).ID

	createStage := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/stage-templates", map[string]any{
		"projectTemplateId": projectID,
		"name":              "Lifecycle Stage",
		"code":              "lifecycle-stage",
		"position":          1,
	}, token)
	if createStage.Code != http.StatusCreated {
		t.Fatalf("create stage failed: %d body=%s", createStage.Code, createStage.Body.String())
	}
	stageID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createStage).ID

	createForm := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/form-templates", map[string]any{
		"stageTemplateId": stageID,
		"name":            "Lifecycle Form",
		"code":            "lifecycle-form",
		"position":        1,
	}, token)
	if createForm.Code != http.StatusCreated {
		t.Fatalf("create form failed: %d body=%s", createForm.Code, createForm.Body.String())
	}
	formID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createForm).ID

	createField := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/field-templates", map[string]any{
		"formTemplateId": formID,
		"name":           "Lifecycle Owner",
		"code":           "lifecycle-owner",
		"position":       1,
		"widgetType":     "input",
	}, token)
	if createField.Code != http.StatusCreated {
		t.Fatalf("create field failed: %d body=%s", createField.Code, createField.Body.String())
	}
	fieldID := decodeJSON[struct {
		ID int `json:"id"`
	}](t, createField).ID

	publishProject := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates/"+itoa(projectID)+"/publish", nil, token)
	if publishProject.Code != http.StatusOK {
		t.Fatalf("publish project failed: %d body=%s", publishProject.Code, publishProject.Body.String())
	}

	publishStage := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/stage-templates/"+itoa(stageID)+"/publish", nil, token)
	if publishStage.Code != http.StatusOK {
		t.Fatalf("publish stage failed: %d body=%s", publishStage.Code, publishStage.Body.String())
	}

	publishForm := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/form-templates/"+itoa(formID)+"/publish", nil, token)
	if publishForm.Code != http.StatusOK {
		t.Fatalf("publish form failed: %d body=%s", publishForm.Code, publishForm.Body.String())
	}

	publishField := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/field-templates/"+itoa(fieldID)+"/publish", nil, token)
	if publishField.Code != http.StatusOK {
		t.Fatalf("publish field failed: %d body=%s", publishField.Code, publishField.Body.String())
	}

	instantiate := doJSONRequest(t, h, http.MethodPost, "/api/tmpl/project-templates/"+itoa(projectID)+"/instantiate", map[string]any{
		"name": "Lifecycle Runtime",
	}, token)
	if instantiate.Code != http.StatusCreated {
		t.Fatalf("instantiate failed: %d body=%s", instantiate.Code, instantiate.Body.String())
	}

	lifecycle := doJSONRequest(t, h, http.MethodGet, "/api/tmpl/project-templates/"+itoa(projectID)+"/lifecycle", nil, token)
	if lifecycle.Code != http.StatusOK {
		t.Fatalf("lifecycle failed: %d body=%s", lifecycle.Code, lifecycle.Body.String())
	}
	payload := decodeJSON[projectTemplateLifecycleResp](t, lifecycle)

	if payload.ProjectTemplate.ID != projectID {
		t.Fatalf("expected lifecycle projectTemplate.id=%d, got %d", projectID, payload.ProjectTemplate.ID)
	}
	if payload.IndustryTemplate.ID != industryID {
		t.Fatalf("expected lifecycle industryTemplate.id=%d, got %d", industryID, payload.IndustryTemplate.ID)
	}
	if payload.Counts.PublishedStages != 1 {
		t.Fatalf("expected publishedStages=1, got %d", payload.Counts.PublishedStages)
	}
	if payload.Counts.PublishedForms != 1 {
		t.Fatalf("expected publishedForms=1, got %d", payload.Counts.PublishedForms)
	}
	if payload.Counts.PublishedFields != 1 {
		t.Fatalf("expected publishedFields=1, got %d", payload.Counts.PublishedFields)
	}
	if payload.Counts.RuntimeProjects != 1 {
		t.Fatalf("expected runtimeProjects=1, got %d", payload.Counts.RuntimeProjects)
	}
	if payload.RuntimeByStageStatus.Active != 1 {
		t.Fatalf("expected runtimeByStageStatus.active=1, got %d", payload.RuntimeByStageStatus.Active)
	}
	if payload.RuntimeByStageStatus.Pending != 0 {
		t.Fatalf("expected runtimeByStageStatus.pending=0, got %d", payload.RuntimeByStageStatus.Pending)
	}
	if payload.RuntimeByStageStatus.Done != 0 {
		t.Fatalf("expected runtimeByStageStatus.done=0, got %d", payload.RuntimeByStageStatus.Done)
	}
	if strings.TrimSpace(payload.Guidance) == "" {
		t.Fatalf("expected lifecycle guidance to be non-empty")
	}
	if payload.Guidance != "Template has active runtime usage. Prefer creating next-version draft and roll out gradually." {
		t.Fatalf("unexpected lifecycle guidance: %s", payload.Guidance)
	}
}
