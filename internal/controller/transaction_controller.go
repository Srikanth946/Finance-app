package controller

import (
	"finance_app/internal/models"
	"finance_app/internal/service"
	"net/http"
	"strconv"

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

func (c *TransactionController) GetDashboardSummary(ctx *gin.Context) {
	summary, err := c.service.GetDashboardSummary()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, summary)
}

func (c *TransactionController) MarkAsPaid(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.service.MarkAsPaid(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Transaction marked as paid"})
}
