package controller

import (
	"finance_app/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	service     service.DashboardService
	interestSvc service.InterestService
}

func NewDashboardController(s service.DashboardService, i service.InterestService) *DashboardController {
	return &DashboardController{service: s, interestSvc: i}
}

func (c *DashboardController) GetSummary(ctx *gin.Context) {
	summary, err := c.service.GetDashboardSummary(c.interestSvc)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, summary)
}
