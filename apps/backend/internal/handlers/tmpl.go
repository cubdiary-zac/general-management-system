package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gms/backend/internal/middleware"
	"gms/backend/internal/models"
)

type TemplateHandler struct {
	db *gorm.DB
}

func NewTemplateHandler(db *gorm.DB) *TemplateHandler {
	return &TemplateHandler{db: db}
}

type createIndustryTemplateRequest struct {
	Name        string                `json:"name"`
	Code        string                `json:"code"`
	Description string                `json:"description"`
	Version     *int                  `json:"version"`
	Status      models.TemplateStatus `json:"status"`
}

type createProjectTemplateRequest struct {
	IndustryTemplateID uint                  `json:"industryTemplateId"`
	Name               string                `json:"name"`
	Code               string                `json:"code"`
	Description        string                `json:"description"`
	Version            *int                  `json:"version"`
	Status             models.TemplateStatus `json:"status"`
}

type createStageTemplateRequest struct {
	ProjectTemplateID uint                  `json:"projectTemplateId"`
	Name              string                `json:"name"`
	Code              string                `json:"code"`
	Description       string                `json:"description"`
	Version           *int                  `json:"version"`
	Status            models.TemplateStatus `json:"status"`
	Position          *int                  `json:"position"`
}

type createFormTemplateRequest struct {
	StageTemplateID uint                  `json:"stageTemplateId"`
	Name            string                `json:"name"`
	Code            string                `json:"code"`
	Description     string                `json:"description"`
	Version         *int                  `json:"version"`
	Status          models.TemplateStatus `json:"status"`
	Position        *int                  `json:"position"`
}

type createFormFieldTemplateRequest struct {
	FormTemplateID uint                       `json:"formTemplateId"`
	Name           string                     `json:"name"`
	Code           string                     `json:"code"`
	Description    string                     `json:"description"`
	Version        *int                       `json:"version"`
	Status         models.TemplateStatus      `json:"status"`
	Position       *int                       `json:"position"`
	WidgetType     models.FormFieldWidgetType `json:"widgetType"`
}

type instantiateProjectTemplateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type instantiateProjectTemplateResponse struct {
	Project models.RuntimeProject        `json:"project"`
	Stages  []models.RuntimeProjectStage `json:"stages"`
	Forms   []models.RuntimeProjectForm  `json:"forms"`
	Fields  []models.RuntimeProjectField `json:"fields"`
}

type instantiateProjectTemplateError struct {
	status  int
	message string
}

func (e *instantiateProjectTemplateError) Error() string {
	return e.message
}

func (h *TemplateHandler) ListIndustryTemplates(c *gin.Context) {
	filters, ok := parseTemplateListCommonFilters(c)
	if !ok {
		return
	}

	items := make([]models.IndustryTemplate, 0)
	query := h.db.Order("id desc")
	if filters.status != nil {
		query = query.Where("status = ?", *filters.status)
	}
	if filters.version != nil {
		query = query.Where("version = ?", *filters.version)
	}

	if err := query.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list industry templates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *TemplateHandler) CreateIndustryTemplate(c *gin.Context) {
	var req createIndustryTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	version, ok := normalizeVersion(req.Version)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "version must be greater than or equal to 1"})
		return
	}

	status, ok := normalizeTemplateStatus(req.Status)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	item := models.IndustryTemplate{
		Name:        name,
		Code:        strings.TrimSpace(req.Code),
		Description: strings.TrimSpace(req.Description),
		Version:     version,
		Status:      status,
	}

	if err := h.db.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create industry template"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *TemplateHandler) ListProjectTemplates(c *gin.Context) {
	filters, ok := parseTemplateListCommonFilters(c)
	if !ok {
		return
	}

	industryTemplateID, ok := parsePositiveIntQuery(c, "industryTemplateId")
	if !ok {
		return
	}

	items := make([]models.ProjectTemplate, 0)
	query := h.db.Order("id desc")
	if filters.status != nil {
		query = query.Where("status = ?", *filters.status)
	}
	if filters.version != nil {
		query = query.Where("version = ?", *filters.version)
	}
	if industryTemplateID != nil {
		query = query.Where("industry_template_id = ?", *industryTemplateID)
	}

	if err := query.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list project templates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *TemplateHandler) CreateProjectTemplate(c *gin.Context) {
	var req createProjectTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if req.IndustryTemplateID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "industryTemplateId is required"})
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	version, ok := normalizeVersion(req.Version)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "version must be greater than or equal to 1"})
		return
	}

	status, ok := normalizeTemplateStatus(req.Status)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	var industryTemplate models.IndustryTemplate
	if err := h.db.Select("id").First(&industryTemplate, req.IndustryTemplateID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "industry template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify industry template"})
		return
	}

	item := models.ProjectTemplate{
		IndustryTemplateID: req.IndustryTemplateID,
		Name:               name,
		Code:               strings.TrimSpace(req.Code),
		Description:        strings.TrimSpace(req.Description),
		Version:            version,
		Status:             status,
	}

	if err := h.db.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create project template"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *TemplateHandler) ListStageTemplates(c *gin.Context) {
	filters, ok := parseTemplateListCommonFilters(c)
	if !ok {
		return
	}

	projectTemplateID, ok := parsePositiveIntQuery(c, "projectTemplateId")
	if !ok {
		return
	}

	items := make([]models.StageTemplate, 0)
	query := h.db.Order("project_template_id asc, position asc, id asc")
	if filters.status != nil {
		query = query.Where("status = ?", *filters.status)
	}
	if filters.version != nil {
		query = query.Where("version = ?", *filters.version)
	}
	if projectTemplateID != nil {
		query = query.Where("project_template_id = ?", *projectTemplateID)
	}

	if err := query.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list stage templates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *TemplateHandler) CreateStageTemplate(c *gin.Context) {
	var req createStageTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if req.ProjectTemplateID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "projectTemplateId is required"})
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	version, ok := normalizeVersion(req.Version)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "version must be greater than or equal to 1"})
		return
	}

	status, ok := normalizeTemplateStatus(req.Status)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	position, ok := normalizePosition(req.Position)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "position must be greater than or equal to 1"})
		return
	}

	var projectTemplate models.ProjectTemplate
	if err := h.db.Select("id").First(&projectTemplate, req.ProjectTemplateID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "project template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify project template"})
		return
	}

	item := models.StageTemplate{
		ProjectTemplateID: req.ProjectTemplateID,
		Name:              name,
		Code:              strings.TrimSpace(req.Code),
		Description:       strings.TrimSpace(req.Description),
		Version:           version,
		Status:            status,
		Position:          position,
	}

	if err := h.db.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create stage template"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *TemplateHandler) ListFormTemplates(c *gin.Context) {
	filters, ok := parseTemplateListCommonFilters(c)
	if !ok {
		return
	}

	stageTemplateID, ok := parsePositiveIntQuery(c, "stageTemplateId")
	if !ok {
		return
	}

	items := make([]models.FormTemplate, 0)
	query := h.db.Order("stage_template_id asc, position asc, id asc")
	if filters.status != nil {
		query = query.Where("status = ?", *filters.status)
	}
	if filters.version != nil {
		query = query.Where("version = ?", *filters.version)
	}
	if stageTemplateID != nil {
		query = query.Where("stage_template_id = ?", *stageTemplateID)
	}

	if err := query.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list form templates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *TemplateHandler) CreateFormTemplate(c *gin.Context) {
	var req createFormTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if req.StageTemplateID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "stageTemplateId is required"})
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	version, ok := normalizeVersion(req.Version)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "version must be greater than or equal to 1"})
		return
	}

	status, ok := normalizeTemplateStatus(req.Status)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	position, ok := normalizePosition(req.Position)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "position must be greater than or equal to 1"})
		return
	}

	var stageTemplate models.StageTemplate
	if err := h.db.Select("id").First(&stageTemplate, req.StageTemplateID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "stage template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify stage template"})
		return
	}

	item := models.FormTemplate{
		StageTemplateID: req.StageTemplateID,
		Name:            name,
		Code:            strings.TrimSpace(req.Code),
		Description:     strings.TrimSpace(req.Description),
		Version:         version,
		Status:          status,
		Position:        position,
	}

	if err := h.db.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create form template"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *TemplateHandler) ListFormFieldTemplates(c *gin.Context) {
	filters, ok := parseTemplateListCommonFilters(c)
	if !ok {
		return
	}

	formTemplateID, ok := parsePositiveIntQuery(c, "formTemplateId")
	if !ok {
		return
	}

	items := make([]models.FormFieldTemplate, 0)
	query := h.db.Order("form_template_id asc, position asc, id asc")
	if filters.status != nil {
		query = query.Where("status = ?", *filters.status)
	}
	if filters.version != nil {
		query = query.Where("version = ?", *filters.version)
	}
	if formTemplateID != nil {
		query = query.Where("form_template_id = ?", *formTemplateID)
	}

	if err := query.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list form field templates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *TemplateHandler) CreateFormFieldTemplate(c *gin.Context) {
	var req createFormFieldTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if req.FormTemplateID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formTemplateId is required"})
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	if req.WidgetType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "widgetType is required"})
		return
	}
	if !req.WidgetType.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid widgetType"})
		return
	}

	version, ok := normalizeVersion(req.Version)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "version must be greater than or equal to 1"})
		return
	}

	status, ok := normalizeTemplateStatus(req.Status)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	position, ok := normalizePosition(req.Position)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "position must be greater than or equal to 1"})
		return
	}

	var formTemplate models.FormTemplate
	if err := h.db.Select("id").First(&formTemplate, req.FormTemplateID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "form template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify form template"})
		return
	}

	item := models.FormFieldTemplate{
		FormTemplateID: req.FormTemplateID,
		Name:           name,
		Code:           strings.TrimSpace(req.Code),
		Description:    strings.TrimSpace(req.Description),
		Version:        version,
		Status:         status,
		Position:       position,
		WidgetType:     req.WidgetType,
	}

	if err := h.db.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create form field template"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *TemplateHandler) InstantiateProjectTemplate(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	templateID, ok := parsePositiveUintParam(c, "id")
	if !ok {
		return
	}

	var req instantiateProjectTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	description := strings.TrimSpace(req.Description)

	resp := instantiateProjectTemplateResponse{
		Stages: make([]models.RuntimeProjectStage, 0),
		Forms:  make([]models.RuntimeProjectForm, 0),
		Fields: make([]models.RuntimeProjectField, 0),
	}

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		var projectTemplate models.ProjectTemplate
		if err := tx.First(&projectTemplate, templateID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &instantiateProjectTemplateError{status: http.StatusNotFound, message: "project template not found"}
			}
			return err
		}

		if projectTemplate.Status != models.TemplateStatusPublished {
			return &instantiateProjectTemplateError{status: http.StatusBadRequest, message: "cannot instantiate project template: template must be published"}
		}

		var industryTemplate models.IndustryTemplate
		if err := tx.Select("id", "status").First(&industryTemplate, projectTemplate.IndustryTemplateID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &instantiateProjectTemplateError{status: http.StatusBadRequest, message: "cannot instantiate project template: parent industry template not found"}
			}
			return err
		}
		if industryTemplate.Status != models.TemplateStatusPublished {
			return &instantiateProjectTemplateError{status: http.StatusBadRequest, message: "cannot instantiate project template: industry template must be published first"}
		}

		stageTemplates := make([]models.StageTemplate, 0)
		if err := tx.Where("project_template_id = ? AND status = ?", projectTemplate.ID, models.TemplateStatusPublished).
			Order("position asc, id asc").
			Find(&stageTemplates).Error; err != nil {
			return err
		}
		if len(stageTemplates) == 0 {
			return &instantiateProjectTemplateError{status: http.StatusBadRequest, message: "cannot instantiate project template: no published stage templates found"}
		}

		formTemplates := make([]models.FormTemplate, 0)
		for _, stageTemplate := range stageTemplates {
			stageForms := make([]models.FormTemplate, 0)
			if err := tx.Where("stage_template_id = ? AND status = ?", stageTemplate.ID, models.TemplateStatusPublished).
				Order("position asc, id asc").
				Find(&stageForms).Error; err != nil {
				return err
			}
			formTemplates = append(formTemplates, stageForms...)
		}

		fieldTemplates := make([]models.FormFieldTemplate, 0)
		for _, formTemplate := range formTemplates {
			formFields := make([]models.FormFieldTemplate, 0)
			if err := tx.Where("form_template_id = ? AND status = ?", formTemplate.ID, models.TemplateStatusPublished).
				Order("position asc, id asc").
				Find(&formFields).Error; err != nil {
				return err
			}
			fieldTemplates = append(fieldTemplates, formFields...)
		}

		project := models.RuntimeProject{
			Name:                   name,
			Description:            description,
			IndustryTemplateID:     projectTemplate.IndustryTemplateID,
			ProjectTemplateID:      projectTemplate.ID,
			ProjectTemplateVersion: projectTemplate.Version,
			CreatedBy:              userID,
		}
		if err := tx.Create(&project).Error; err != nil {
			return err
		}
		resp.Project = project

		stageTemplateIDToRuntimeStageID := make(map[uint]uint, len(stageTemplates))
		for idx, stageTemplate := range stageTemplates {
			stageStatus := models.RuntimeProjectStageStatusPending
			if idx == 0 {
				stageStatus = models.RuntimeProjectStageStatusActive
			}

			stage := models.RuntimeProjectStage{
				RuntimeProjectID: project.ID,
				StageTemplateID:  stageTemplate.ID,
				Name:             stageTemplate.Name,
				Code:             stageTemplate.Code,
				Description:      stageTemplate.Description,
				Position:         stageTemplate.Position,
				Status:           stageStatus,
			}
			if err := tx.Create(&stage).Error; err != nil {
				return err
			}
			resp.Stages = append(resp.Stages, stage)
			stageTemplateIDToRuntimeStageID[stageTemplate.ID] = stage.ID
		}

		formTemplateIDToRuntimeFormID := make(map[uint]uint, len(formTemplates))
		for _, formTemplate := range formTemplates {
			runtimeStageID, exists := stageTemplateIDToRuntimeStageID[formTemplate.StageTemplateID]
			if !exists {
				return errors.New("failed to map stage template to runtime stage")
			}

			form := models.RuntimeProjectForm{
				RuntimeProjectID:      project.ID,
				RuntimeProjectStageID: runtimeStageID,
				FormTemplateID:        formTemplate.ID,
				Name:                  formTemplate.Name,
				Code:                  formTemplate.Code,
				Description:           formTemplate.Description,
				Position:              formTemplate.Position,
			}
			if err := tx.Create(&form).Error; err != nil {
				return err
			}
			resp.Forms = append(resp.Forms, form)
			formTemplateIDToRuntimeFormID[formTemplate.ID] = form.ID
		}

		for _, fieldTemplate := range fieldTemplates {
			runtimeFormID, exists := formTemplateIDToRuntimeFormID[fieldTemplate.FormTemplateID]
			if !exists {
				return errors.New("failed to map form template to runtime form")
			}

			field := models.RuntimeProjectField{
				RuntimeProjectID:     project.ID,
				RuntimeProjectFormID: runtimeFormID,
				FormFieldTemplateID:  fieldTemplate.ID,
				Name:                 fieldTemplate.Name,
				Code:                 fieldTemplate.Code,
				Description:          fieldTemplate.Description,
				Position:             fieldTemplate.Position,
				WidgetType:           fieldTemplate.WidgetType,
				ValueText:            nil,
			}
			if err := tx.Create(&field).Error; err != nil {
				return err
			}
			resp.Fields = append(resp.Fields, field)
		}

		return nil
	}); err != nil {
		var instantiateErr *instantiateProjectTemplateError
		if errors.As(err, &instantiateErr) {
			c.JSON(instantiateErr.status, gin.H{"error": instantiateErr.message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to instantiate project template"})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *TemplateHandler) PublishIndustryTemplate(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	templateID, ok := parsePositiveUintParam(c, "id")
	if !ok {
		return
	}

	var item models.IndustryTemplate
	if err := h.db.First(&item, templateID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "industry template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load industry template"})
		return
	}

	publishedAt := time.Now().UTC()
	publishedBy := userID
	item.Status = models.TemplateStatusPublished
	item.PublishedAt = &publishedAt
	item.PublishedBy = &publishedBy

	if err := h.db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish industry template"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *TemplateHandler) UnpublishIndustryTemplate(c *gin.Context) {
	templateID, ok := parsePositiveUintParam(c, "id")
	if !ok {
		return
	}

	var item models.IndustryTemplate
	if err := h.db.First(&item, templateID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "industry template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load industry template"})
		return
	}

	var publishedChildren int64
	if err := h.db.Model(&models.ProjectTemplate{}).
		Where("industry_template_id = ? AND status = ?", item.ID, models.TemplateStatusPublished).
		Count(&publishedChildren).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify project templates"})
		return
	}
	if publishedChildren > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot unpublish industry template: published project templates exist"})
		return
	}

	item.Status = models.TemplateStatusDraft
	item.PublishedAt = nil
	item.PublishedBy = nil

	if err := h.db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unpublish industry template"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *TemplateHandler) PublishProjectTemplate(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	templateID, ok := parsePositiveUintParam(c, "id")
	if !ok {
		return
	}

	var item models.ProjectTemplate
	if err := h.db.First(&item, templateID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "project template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load project template"})
		return
	}

	var parent models.IndustryTemplate
	if err := h.db.Select("id", "status").First(&parent, item.IndustryTemplateID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot publish project template: parent industry template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify parent industry template"})
		return
	}
	if parent.Status != models.TemplateStatusPublished {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot publish project template: industry template must be published first"})
		return
	}

	publishedAt := time.Now().UTC()
	publishedBy := userID
	item.Status = models.TemplateStatusPublished
	item.PublishedAt = &publishedAt
	item.PublishedBy = &publishedBy

	if err := h.db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish project template"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *TemplateHandler) UnpublishProjectTemplate(c *gin.Context) {
	templateID, ok := parsePositiveUintParam(c, "id")
	if !ok {
		return
	}

	var item models.ProjectTemplate
	if err := h.db.First(&item, templateID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "project template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load project template"})
		return
	}

	var publishedChildren int64
	if err := h.db.Model(&models.StageTemplate{}).
		Where("project_template_id = ? AND status = ?", item.ID, models.TemplateStatusPublished).
		Count(&publishedChildren).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify stage templates"})
		return
	}
	if publishedChildren > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot unpublish project template: published stage templates exist"})
		return
	}

	item.Status = models.TemplateStatusDraft
	item.PublishedAt = nil
	item.PublishedBy = nil

	if err := h.db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unpublish project template"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *TemplateHandler) PublishStageTemplate(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	templateID, ok := parsePositiveUintParam(c, "id")
	if !ok {
		return
	}

	var item models.StageTemplate
	if err := h.db.First(&item, templateID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "stage template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load stage template"})
		return
	}

	var parent models.ProjectTemplate
	if err := h.db.Select("id", "status").First(&parent, item.ProjectTemplateID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot publish stage template: parent project template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify parent project template"})
		return
	}
	if parent.Status != models.TemplateStatusPublished {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot publish stage template: project template must be published first"})
		return
	}

	publishedAt := time.Now().UTC()
	publishedBy := userID
	item.Status = models.TemplateStatusPublished
	item.PublishedAt = &publishedAt
	item.PublishedBy = &publishedBy

	if err := h.db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish stage template"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *TemplateHandler) UnpublishStageTemplate(c *gin.Context) {
	templateID, ok := parsePositiveUintParam(c, "id")
	if !ok {
		return
	}

	var item models.StageTemplate
	if err := h.db.First(&item, templateID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "stage template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load stage template"})
		return
	}

	var publishedChildren int64
	if err := h.db.Model(&models.FormTemplate{}).
		Where("stage_template_id = ? AND status = ?", item.ID, models.TemplateStatusPublished).
		Count(&publishedChildren).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify form templates"})
		return
	}
	if publishedChildren > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot unpublish stage template: published form templates exist"})
		return
	}

	item.Status = models.TemplateStatusDraft
	item.PublishedAt = nil
	item.PublishedBy = nil

	if err := h.db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unpublish stage template"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *TemplateHandler) PublishFormTemplate(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	templateID, ok := parsePositiveUintParam(c, "id")
	if !ok {
		return
	}

	var item models.FormTemplate
	if err := h.db.First(&item, templateID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "form template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load form template"})
		return
	}

	var parent models.StageTemplate
	if err := h.db.Select("id", "status").First(&parent, item.StageTemplateID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot publish form template: parent stage template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify parent stage template"})
		return
	}
	if parent.Status != models.TemplateStatusPublished {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot publish form template: stage template must be published first"})
		return
	}

	publishedAt := time.Now().UTC()
	publishedBy := userID
	item.Status = models.TemplateStatusPublished
	item.PublishedAt = &publishedAt
	item.PublishedBy = &publishedBy

	if err := h.db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish form template"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *TemplateHandler) UnpublishFormTemplate(c *gin.Context) {
	templateID, ok := parsePositiveUintParam(c, "id")
	if !ok {
		return
	}

	var item models.FormTemplate
	if err := h.db.First(&item, templateID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "form template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load form template"})
		return
	}

	var publishedChildren int64
	if err := h.db.Model(&models.FormFieldTemplate{}).
		Where("form_template_id = ? AND status = ?", item.ID, models.TemplateStatusPublished).
		Count(&publishedChildren).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify field templates"})
		return
	}
	if publishedChildren > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot unpublish form template: published field templates exist"})
		return
	}

	item.Status = models.TemplateStatusDraft
	item.PublishedAt = nil
	item.PublishedBy = nil

	if err := h.db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unpublish form template"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *TemplateHandler) PublishFormFieldTemplate(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	templateID, ok := parsePositiveUintParam(c, "id")
	if !ok {
		return
	}

	var item models.FormFieldTemplate
	if err := h.db.First(&item, templateID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "field template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load field template"})
		return
	}

	var parent models.FormTemplate
	if err := h.db.Select("id", "status").First(&parent, item.FormTemplateID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot publish field template: parent form template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify parent form template"})
		return
	}
	if parent.Status != models.TemplateStatusPublished {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot publish field template: form template must be published first"})
		return
	}

	publishedAt := time.Now().UTC()
	publishedBy := userID
	item.Status = models.TemplateStatusPublished
	item.PublishedAt = &publishedAt
	item.PublishedBy = &publishedBy

	if err := h.db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish field template"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *TemplateHandler) UnpublishFormFieldTemplate(c *gin.Context) {
	templateID, ok := parsePositiveUintParam(c, "id")
	if !ok {
		return
	}

	var item models.FormFieldTemplate
	if err := h.db.First(&item, templateID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "field template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load field template"})
		return
	}

	item.Status = models.TemplateStatusDraft
	item.PublishedAt = nil
	item.PublishedBy = nil

	if err := h.db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unpublish field template"})
		return
	}

	c.JSON(http.StatusOK, item)
}

type templateListCommonFilters struct {
	status  *models.TemplateStatus
	version *int
}

func parseTemplateListCommonFilters(c *gin.Context) (templateListCommonFilters, bool) {
	filters := templateListCommonFilters{}

	statusParam := strings.TrimSpace(c.Query("status"))
	if statusParam != "" {
		status := models.TemplateStatus(statusParam)
		if !status.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
			return filters, false
		}
		filters.status = &status
	}

	version, ok := parsePositiveIntQuery(c, "version")
	if !ok {
		return filters, false
	}
	filters.version = version

	return filters, true
}

func parsePositiveIntQuery(c *gin.Context, key string) (*int, bool) {
	value := strings.TrimSpace(c.Query(key))
	if value == "" {
		return nil, true
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": key + " must be a positive integer"})
		return nil, false
	}

	return &parsed, true
}

func parsePositiveUintParam(c *gin.Context, key string) (uint, bool) {
	value := strings.TrimSpace(c.Param(key))
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": key + " must be a positive integer"})
		return 0, false
	}
	return uint(parsed), true
}

func currentUserID(c *gin.Context) (uint, bool) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return 0, false
	}
	return user.ID, true
}

func normalizeTemplateStatus(status models.TemplateStatus) (models.TemplateStatus, bool) {
	if status == "" {
		return models.TemplateStatusDraft, true
	}

	if !status.IsValid() {
		return "", false
	}

	return status, true
}

func normalizeVersion(version *int) (int, bool) {
	if version == nil {
		return 1, true
	}
	if *version < 1 {
		return 0, false
	}
	return *version, true
}

func normalizePosition(position *int) (int, bool) {
	if position == nil {
		return 1, true
	}
	if *position < 1 {
		return 0, false
	}
	return *position, true
}
