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

type CRMHandler struct {
	db *gorm.DB
}

func NewCRMHandler(db *gorm.DB) *CRMHandler {
	return &CRMHandler{db: db}
}

type createCustomerRequest struct {
	Name    string `json:"name"`
	Company string `json:"company"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
}

type createLeadRequest struct {
	CustomerID *uint    `json:"customerId"`
	Name       string   `json:"name"`
	Source     string   `json:"source"`
	Amount     *float64 `json:"amount"`
}

type patchLeadStatusRequest struct {
	Status models.LeadStatus `json:"status"`
}

type leadSummaryRow struct {
	Status models.LeadStatus
	Total  int64
}

func (h *CRMHandler) ListCustomers(c *gin.Context) {
	customers := make([]models.Customer, 0)
	if err := h.db.Order("id desc").Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list customers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": customers})
}

func (h *CRMHandler) CreateCustomer(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req createCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer name is required"})
		return
	}

	customer := models.Customer{
		Name:    req.Name,
		Company: strings.TrimSpace(req.Company),
		Email:   strings.TrimSpace(req.Email),
		Phone:   strings.TrimSpace(req.Phone),
		OwnerID: user.ID,
	}

	if err := h.db.Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create customer"})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

func (h *CRMHandler) ListLeads(c *gin.Context) {
	statusParam := strings.TrimSpace(c.Query("status"))
	keyword := strings.TrimSpace(c.Query("q"))
	leads := make([]models.Lead, 0)

	query := h.db.Order("id desc")
	if statusParam != "" {
		status := models.LeadStatus(statusParam)
		if !status.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
			return
		}
		query = query.Where("status = ?", status)
	}

	if keyword != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(keyword)+"%")
	}

	if err := query.Find(&leads).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list leads"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": leads})
}

func (h *CRMHandler) CreateLead(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req createLeadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lead name is required"})
		return
	}

	req.Source = strings.TrimSpace(req.Source)
	if req.Source == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lead source is required"})
		return
	}

	if req.CustomerID != nil {
		if *req.CustomerID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "customerId must be a positive integer"})
			return
		}

		var customer models.Customer
		if err := h.db.Select("id").First(&customer, *req.CustomerID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify customer"})
			return
		}
	}

	if req.Amount != nil && *req.Amount < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be non-negative"})
		return
	}

	lead := models.Lead{
		CustomerID: req.CustomerID,
		Name:       req.Name,
		Source:     req.Source,
		Status:     models.LeadNew,
		Amount:     req.Amount,
		OwnerID:    user.ID,
	}

	if err := h.db.Create(&lead).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create lead"})
		return
	}

	c.JSON(http.StatusCreated, lead)
}

func (h *CRMHandler) PatchLeadStatus(c *gin.Context) {
	if _, ok := middleware.CurrentUser(c); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	leadID, err := strconv.Atoi(c.Param("id"))
	if err != nil || leadID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lead id"})
		return
	}

	var req patchLeadStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if !req.Status.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	var lead models.Lead
	if err := h.db.First(&lead, leadID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "lead not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load lead"})
		return
	}

	if !models.CanTransitionLeadStatus(lead.Status, req.Status) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "invalid status transition",
			"from":        lead.Status,
			"to":          req.Status,
			"allowedNext": allowedNextLeadStatus(lead.Status),
		})
		return
	}

	lead.Status = req.Status
	if err := h.db.Save(&lead).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update lead status"})
		return
	}

	c.JSON(http.StatusOK, lead)
}

func (h *CRMHandler) Summary(c *gin.Context) {
	rows := make([]leadSummaryRow, 0)
	if err := h.db.Model(&models.Lead{}).Select("status, COUNT(*) as total").Group("status").Scan(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to summarize leads"})
		return
	}

	counts := map[models.LeadStatus]int64{
		models.LeadNew:       0,
		models.LeadContacted: 0,
		models.LeadQualified: 0,
		models.LeadWon:       0,
		models.LeadLost:      0,
	}

	for _, row := range rows {
		if row.Status.IsValid() {
			counts[row.Status] = row.Total
		}
	}

	c.JSON(http.StatusOK, gin.H{"counts": counts})
}

func allowedNextLeadStatus(status models.LeadStatus) []models.LeadStatus {
	switch status {
	case models.LeadNew:
		return []models.LeadStatus{models.LeadNew, models.LeadContacted, models.LeadLost}
	case models.LeadContacted:
		return []models.LeadStatus{models.LeadContacted, models.LeadQualified, models.LeadLost}
	case models.LeadQualified:
		return []models.LeadStatus{models.LeadQualified, models.LeadWon, models.LeadLost}
	case models.LeadWon:
		return []models.LeadStatus{models.LeadWon}
	case models.LeadLost:
		return []models.LeadStatus{models.LeadLost}
	default:
		return []models.LeadStatus{}
	}
}
