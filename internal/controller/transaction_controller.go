package controller

import (
	"finance_app/internal/models"
	"finance_app/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	service service.TransactionService
}

func NewTransactionController(s service.TransactionService) *TransactionController {
	return &TransactionController{service: s}
}

func (c *TransactionController) CreateTransaction(ctx *gin.Context) {
	var txn models.Transaction
	if err := ctx.ShouldBindJSON(&txn); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.CreateTransaction(&txn); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, txn)
}

func (c *TransactionController) GetAllTransactions(ctx *gin.Context) {
	txns, err := c.service.GetAllTransactions()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, txns)
}

func (c *TransactionController) MarkAsPaid(ctx *gin.Context) {
	mobileNumber := ctx.Param("id")
	if mobileNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Mobile number is required"})
		return
	}

	if err := c.service.MarkAsPaid(mobileNumber); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Transaction marked as paid"})
}

func (c *TransactionController) ExtendLoan(ctx *gin.Context) {
	mobileNumber := ctx.Param("id")
	if mobileNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Mobile number is required"})
		return
	}

	var req struct {
		InterestRate           float64 `json:"interest_rate"`
		CompoundDurationMonths int     `json:"compound_duration_months"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Interest rate and compound duration are required"})
		return
	}

	if req.InterestRate <= 0 || req.CompoundDurationMonths <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Interest rate and compound duration must be greater than zero"})
		return
	}

	if err := c.service.ExtendLoan(mobileNumber, req.InterestRate, req.CompoundDurationMonths); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Loan extended successfully with new terms"})
}

func (c *TransactionController) GetTransaction(ctx *gin.Context) {
	mobileNumber := ctx.Param("id")
	if mobileNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Mobile number is required"})
		return
	}

	txn, summary, err := c.service.GetLoanWithInterest(mobileNumber)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found or error calculating interest"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"loan_details":     txn,
		"current_interest": summary,
	})
}
