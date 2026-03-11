package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gms/backend/internal/models"
)

type templateGovernanceError struct {
	status  int
	message string
}

func (e *templateGovernanceError) Error() string {
	return e.message
}

type projectTemplateNextVersionResponse struct {
	ProjectTemplate models.ProjectTemplate     `json:"projectTemplate"`
	Stages          []models.StageTemplate     `json:"stages"`
	Forms           []models.FormTemplate      `json:"forms"`
	Fields          []models.FormFieldTemplate `json:"fields"`
}

type projectTemplateLifecycleIndustry struct {
	ID      uint                  `json:"id"`
	Name    string                `json:"name"`
	Code    string                `json:"code"`
	Status  models.TemplateStatus `json:"status"`
	Version int                   `json:"version"`
}

type projectTemplateLifecycleCounts struct {
	PublishedStages int64 `json:"publishedStages"`
	PublishedForms  int64 `json:"publishedForms"`
	PublishedFields int64 `json:"publishedFields"`
	RuntimeProjects int64 `json:"runtimeProjects"`
}

type projectTemplateRuntimeByStageStatus struct {
	Pending int64 `json:"pending"`
	Active  int64 `json:"active"`
	Done    int64 `json:"done"`
}

type projectTemplateLifecycleResponse struct {
	ProjectTemplate      models.ProjectTemplate              `json:"projectTemplate"`
	IndustryTemplate     projectTemplateLifecycleIndustry    `json:"industryTemplate"`
	Counts               projectTemplateLifecycleCounts      `json:"counts"`
	RuntimeByStageStatus projectTemplateRuntimeByStageStatus `json:"runtimeByStageStatus"`
	Guidance             string                              `json:"guidance"`
}

func (h *TemplateHandler) CreateNextProjectTemplateVersion(c *gin.Context) {
	sourceID, ok := parsePositiveUintParam(c, "id")
	if !ok {
		return
	}

	resp := projectTemplateNextVersionResponse{
		Stages: make([]models.StageTemplate, 0),
		Forms:  make([]models.FormTemplate, 0),
		Fields: make([]models.FormFieldTemplate, 0),
	}

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		var source models.ProjectTemplate
		if err := tx.First(&source, sourceID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &templateGovernanceError{status: http.StatusNotFound, message: "project template not found"}
			}
			return err
		}

		if source.Status != models.TemplateStatusPublished {
			return &templateGovernanceError{
				status:  http.StatusBadRequest,
				message: "cannot create next version: source project template must be published",
			}
		}

		nextVersion := source.Version + 1
		nextName, err := resolveNextProjectTemplateName(tx, source.IndustryTemplateID, source.Name, nextVersion)
		if err != nil {
			return err
		}

		nextTemplate := models.ProjectTemplate{
			IndustryTemplateID: source.IndustryTemplateID,
			Name:               nextName,
			Code:               source.Code,
			Description:        source.Description,
			Version:            nextVersion,
			Status:             models.TemplateStatusDraft,
			PublishedAt:        nil,
			PublishedBy:        nil,
		}
		if err := tx.Create(&nextTemplate).Error; err != nil {
			return err
		}
		resp.ProjectTemplate = nextTemplate

		sourceStages := make([]models.StageTemplate, 0)
		if err := tx.Where("project_template_id = ?", source.ID).
			Order("position asc, id asc").
			Find(&sourceStages).Error; err != nil {
			return err
		}

		stageIDMap := make(map[uint]uint, len(sourceStages))
		for _, sourceStage := range sourceStages {
			clonedStage := models.StageTemplate{
				ProjectTemplateID: nextTemplate.ID,
				Name:              sourceStage.Name,
				Code:              sourceStage.Code,
				Description:       sourceStage.Description,
				Version:           sourceStage.Version,
				Status:            models.TemplateStatusDraft,
				PublishedAt:       nil,
				PublishedBy:       nil,
				Position:          sourceStage.Position,
			}
			if err := tx.Create(&clonedStage).Error; err != nil {
				return err
			}
			stageIDMap[sourceStage.ID] = clonedStage.ID
			resp.Stages = append(resp.Stages, clonedStage)

			sourceForms := make([]models.FormTemplate, 0)
			if err := tx.Where("stage_template_id = ?", sourceStage.ID).
				Order("position asc, id asc").
				Find(&sourceForms).Error; err != nil {
				return err
			}

			formIDMap := make(map[uint]uint, len(sourceForms))
			for _, sourceForm := range sourceForms {
				clonedForm := models.FormTemplate{
					StageTemplateID: stageIDMap[sourceStage.ID],
					Name:            sourceForm.Name,
					Code:            sourceForm.Code,
					Description:     sourceForm.Description,
					Version:         sourceForm.Version,
					Status:          models.TemplateStatusDraft,
					PublishedAt:     nil,
					PublishedBy:     nil,
					Position:        sourceForm.Position,
				}
				if err := tx.Create(&clonedForm).Error; err != nil {
					return err
				}
				formIDMap[sourceForm.ID] = clonedForm.ID
				resp.Forms = append(resp.Forms, clonedForm)

				sourceFields := make([]models.FormFieldTemplate, 0)
				if err := tx.Where("form_template_id = ?", sourceForm.ID).
					Order("position asc, id asc").
					Find(&sourceFields).Error; err != nil {
					return err
				}

				for _, sourceField := range sourceFields {
					clonedField := models.FormFieldTemplate{
						FormTemplateID: formIDMap[sourceForm.ID],
						Name:           sourceField.Name,
						Code:           sourceField.Code,
						Description:    sourceField.Description,
						Version:        sourceField.Version,
						Status:         models.TemplateStatusDraft,
						PublishedAt:    nil,
						PublishedBy:    nil,
						Position:       sourceField.Position,
						WidgetType:     sourceField.WidgetType,
					}
					if err := tx.Create(&clonedField).Error; err != nil {
						return err
					}
					resp.Fields = append(resp.Fields, clonedField)
				}
			}
		}

		return nil
	}); err != nil {
		var governanceErr *templateGovernanceError
		if errors.As(err, &governanceErr) {
			c.JSON(governanceErr.status, gin.H{"error": governanceErr.message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create next project template version"})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *TemplateHandler) GetProjectTemplateLifecycle(c *gin.Context) {
	templateID, ok := parsePositiveUintParam(c, "id")
	if !ok {
		return
	}

	var projectTemplate models.ProjectTemplate
	if err := h.db.First(&projectTemplate, templateID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "project template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load project template"})
		return
	}

	var industryTemplate models.IndustryTemplate
	if err := h.db.Select("id", "name", "code", "status", "version").
		First(&industryTemplate, projectTemplate.IndustryTemplateID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "industry template not found for project template"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load industry template"})
		return
	}

	counts := projectTemplateLifecycleCounts{}
	if err := h.db.Model(&models.StageTemplate{}).
		Where("project_template_id = ? AND status = ?", projectTemplate.ID, models.TemplateStatusPublished).
		Count(&counts.PublishedStages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count published stage templates"})
		return
	}

	if err := h.db.Model(&models.FormTemplate{}).
		Joins("JOIN stage_templates ON stage_templates.id = form_templates.stage_template_id").
		Where("stage_templates.project_template_id = ? AND form_templates.status = ?", projectTemplate.ID, models.TemplateStatusPublished).
		Count(&counts.PublishedForms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count published form templates"})
		return
	}

	if err := h.db.Model(&models.FormFieldTemplate{}).
		Joins("JOIN form_templates ON form_templates.id = form_field_templates.form_template_id").
		Joins("JOIN stage_templates ON stage_templates.id = form_templates.stage_template_id").
		Where("stage_templates.project_template_id = ? AND form_field_templates.status = ?", projectTemplate.ID, models.TemplateStatusPublished).
		Count(&counts.PublishedFields).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count published field templates"})
		return
	}

	if err := h.db.Model(&models.RuntimeProject{}).
		Where("project_template_id = ?", projectTemplate.ID).
		Count(&counts.RuntimeProjects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count runtime projects"})
		return
	}

	runtimeByStageStatus := projectTemplateRuntimeByStageStatus{}
	statusCounts := make([]struct {
		Status models.RuntimeProjectStageStatus `json:"status"`
		Count  int64                            `json:"count"`
	}, 0)

	if err := h.db.Model(&models.RuntimeProjectStage{}).
		Select("runtime_project_stages.status AS status, COUNT(*) AS count").
		Joins("JOIN runtime_projects ON runtime_projects.id = runtime_project_stages.runtime_project_id").
		Where("runtime_projects.project_template_id = ?", projectTemplate.ID).
		Group("runtime_project_stages.status").
		Scan(&statusCounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count runtime stage statuses"})
		return
	}

	for _, item := range statusCounts {
		switch item.Status {
		case models.RuntimeProjectStageStatusPending:
			runtimeByStageStatus.Pending = item.Count
		case models.RuntimeProjectStageStatusActive:
			runtimeByStageStatus.Active = item.Count
		case models.RuntimeProjectStageStatusDone:
			runtimeByStageStatus.Done = item.Count
		}
	}

	guidance := "Draft template with runtime history detected. Validate compatibility before publishing."
	if counts.RuntimeProjects == 0 {
		guidance = "No active runtime projects yet. Safe to iterate quickly."
	} else if projectTemplate.Status == models.TemplateStatusPublished {
		guidance = "Template has active runtime usage. Prefer creating next-version draft and roll out gradually."
	}

	resp := projectTemplateLifecycleResponse{
		ProjectTemplate: projectTemplate,
		IndustryTemplate: projectTemplateLifecycleIndustry{
			ID:      industryTemplate.ID,
			Name:    industryTemplate.Name,
			Code:    industryTemplate.Code,
			Status:  industryTemplate.Status,
			Version: industryTemplate.Version,
		},
		Counts:               counts,
		RuntimeByStageStatus: runtimeByStageStatus,
		Guidance:             guidance,
	}

	c.JSON(http.StatusOK, resp)
}

func resolveNextProjectTemplateName(tx *gorm.DB, industryTemplateID uint, baseName string, version int) (string, error) {
	candidate := baseName
	suffix := 2
	for {
		var existingCount int64
		if err := tx.Model(&models.ProjectTemplate{}).
			Where("industry_template_id = ? AND version = ? AND name = ?", industryTemplateID, version, candidate).
			Count(&existingCount).Error; err != nil {
			return "", err
		}
		if existingCount == 0 {
			return candidate, nil
		}

		candidate = fmt.Sprintf("%s (v%d)", baseName, suffix)
		suffix++
	}
}
