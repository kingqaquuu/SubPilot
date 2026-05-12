package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/config"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/response"
)

type HealthHandler struct {
	cfg *config.Config
}

func NewHealthHandler(cfg *config.Config) *HealthHandler {
	return &HealthHandler{cfg: cfg}
}

func (h *HealthHandler) Check(c *gin.Context) {
	response.Success(c, gin.H{
		"service": h.cfg.App.Name,
		"status":  "ok",
	})
}
