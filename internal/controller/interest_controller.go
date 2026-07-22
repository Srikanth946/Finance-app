package controller

import (
	"finance_app/internal/models"
	"finance_app/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InterestController struct {
	service service.InterestService
}

func NewInterestController(s service.InterestService) *InterestController {
	return &InterestController{service: s}
}

func (c *InterestController) Calculate(ctx *gin.Context) {
	// In a real app, we'd fetch the transaction from a service first
	// For now, we expect the transaction object in the request body
	var txn models.Transaction
	if err := ctx.ShouldBindJSON(&txn); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	summary, err := c.service.CalculateInterest(&txn)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, summary)
}
