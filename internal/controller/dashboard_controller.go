package controller

import (
	"finance_app/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	service service.DashboardService
}

func NewDashboardController(s service.DashboardService) *DashboardController {
	return &DashboardController{service: s}
}

func (c *DashboardController) GetSummary(ctx *gin.Context) {
	summary, err := c.service.GetDashboardSummary()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, summary)
}
