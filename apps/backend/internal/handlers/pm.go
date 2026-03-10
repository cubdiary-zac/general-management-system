package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gms/backend/internal/middleware"
	"gms/backend/internal/models"
)

type PMHandler struct {
	db *gorm.DB
}

func NewPMHandler(db *gorm.DB) *PMHandler {
	return &PMHandler{db: db}
}

type createProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type createTaskRequest struct {
	ProjectID   uint   `json:"projectId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	AssigneeID  *uint  `json:"assigneeId"`
}

type patchTaskStatusRequest struct {
	Status models.TaskStatus `json:"status"`
}

func (h *PMHandler) ListProjects(c *gin.Context) {
	projects := make([]models.Project, 0)
	if err := h.db.Order("id desc").Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list projects"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": projects})
}

func (h *PMHandler) CreateProject(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req createProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project name is required"})
		return
	}

	project := models.Project{
		Name:        req.Name,
		Description: strings.TrimSpace(req.Description),
		OwnerID:     user.ID,
	}

	if err := h.db.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, project)
}

func (h *PMHandler) ListTasks(c *gin.Context) {
	projectIDParam := strings.TrimSpace(c.Query("projectId"))
	tasks := make([]models.Task, 0)

	query := h.db.Order("id desc")
	if projectIDParam != "" {
		projectID, err := strconv.Atoi(projectIDParam)
		if err != nil || projectID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "projectId must be a positive integer"})
			return
		}
		query = query.Where("project_id = ?", projectID)
	}

	if err := query.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": tasks})
}

func (h *PMHandler) CreateTask(c *gin.Context) {
	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if req.ProjectID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "projectId is required"})
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	var project models.Project
	if err := h.db.First(&project, req.ProjectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify project"})
		return
	}

	task := models.Task{
		ProjectID:   req.ProjectID,
		Title:       req.Title,
		Description: strings.TrimSpace(req.Description),
		AssigneeID:  req.AssigneeID,
		Status:      models.TaskTodo,
	}

	if err := h.db.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *PMHandler) PatchTaskStatus(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	var req patchTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if !req.Status.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	var task models.Task
	if err := h.db.First(&task, taskID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load task"})
		return
	}

	if !models.CanTransitionTaskStatus(task.Status, req.Status) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "invalid status transition",
			"from":        task.Status,
			"to":          req.Status,
			"allowedNext": allowedNextStatus(task.Status),
		})
		return
	}

	task.Status = req.Status
	if err := h.db.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task status"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func allowedNextStatus(status models.TaskStatus) []models.TaskStatus {
	switch status {
	case models.TaskTodo:
		return []models.TaskStatus{models.TaskTodo, models.TaskInProgress}
	case models.TaskInProgress:
		return []models.TaskStatus{models.TaskInProgress, models.TaskInReview}
	case models.TaskInReview:
		return []models.TaskStatus{models.TaskInReview, models.TaskDone}
	case models.TaskDone:
		return []models.TaskStatus{models.TaskDone}
	default:
		return []models.TaskStatus{}
	}
}
