package router

import (
	"finance_app/internal/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	txnController *controller.TransactionController,
	dashController *controller.DashboardController,
	intController *controller.InterestController,
) *gin.Engine {
	r := gin.Default()

	// Transaction Routes
	r.POST("/transactions", txnController.CreateTransaction)
	r.GET("/transactions", txnController.GetAllTransactions)
	r.POST("/transactions/paid", txnController.MarkAsPaid)

	// Dashboard Routes
	r.GET("/dashboard", dashController.GetSummary)

	// Interest Routes
	r.POST("/interest/calculate", intController.Calculate)

	return r
}
