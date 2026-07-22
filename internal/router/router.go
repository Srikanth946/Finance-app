package router

import (
	"finance_app/internal/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter(txnController *controller.TransactionController) *gin.Engine {
	r := gin.Default()

	r.POST("/transactions", txnController.CreateTransaction)
	r.GET("/transactions", txnController.GetAllTransactions)
	r.POST("/transactions/paid", txnController.MarkAsPaid)
	r.GET("/dashboard", txnController.GetDashboardSummary)

	return r
}
