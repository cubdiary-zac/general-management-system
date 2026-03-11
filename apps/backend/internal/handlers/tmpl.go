package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

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

func (h *TemplateHandler) ListIndustryTemplates(c *gin.Context) {
	items := make([]models.IndustryTemplate, 0)
	if err := h.db.Order("id desc").Find(&items).Error; err != nil {
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
	items := make([]models.ProjectTemplate, 0)
	if err := h.db.Order("id desc").Find(&items).Error; err != nil {
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
	items := make([]models.StageTemplate, 0)
	if err := h.db.Order("project_template_id asc, position asc, id asc").Find(&items).Error; err != nil {
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
	items := make([]models.FormTemplate, 0)
	if err := h.db.Order("stage_template_id asc, position asc, id asc").Find(&items).Error; err != nil {
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
	items := make([]models.FormFieldTemplate, 0)
	if err := h.db.Order("form_template_id asc, position asc, id asc").Find(&items).Error; err != nil {
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
