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
	// Check if the request is a general calculation (amount, rate, months)
	var generalReq struct {
		Amount           float64 `json:"amount"`
		InterestRate     float64 `json:"interest_rate"`
		Months           int     `json:"months"`
		CompoundDuration int     `json:"compound_duration_months"`
	}

	if err := ctx.ShouldBindJSON(&generalReq); err == nil && generalReq.Amount != 0 {
		summary, err := c.service.CalculateGeneralInterest(
			generalReq.Amount,
			generalReq.InterestRate,
			generalReq.Months,
			generalReq.CompoundDuration,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, summary)
		return
	}

	// Fallback to transaction-based calculation
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
