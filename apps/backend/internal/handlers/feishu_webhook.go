package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gms/backend/internal/config"
)

type FeishuWebhookHandler struct {
	verificationToken string
}

func NewFeishuWebhookHandler(cfg config.Config) *FeishuWebhookHandler {
	return &FeishuWebhookHandler{
		verificationToken: cfg.FeishuVerificationToken,
	}
}

type feishuCallbackRequest struct {
	Type      string `json:"type"`
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
}

func (h *FeishuWebhookHandler) Callback(c *gin.Context) {
	var req feishuCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40001,
			"msg":  "invalid json body",
		})
		return
	}

	if h.verificationToken != "" && req.Token != h.verificationToken {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 40003,
			"msg":  "invalid token",
		})
		return
	}

	// Feishu URL verification handshake
	if req.Type == "url_verification" {
		c.JSON(http.StatusOK, gin.H{
			"challenge": req.Challenge,
		})
		return
	}

	// Minimal viable implementation for event callback acknowledgement.
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
	})
}
