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
	if txnController != nil {
		r.POST("/transactions", txnController.CreateTransaction)
		r.GET("/transactions", txnController.GetAllTransactions)
		r.GET("/transactions/:id", txnController.GetTransaction)
		r.PUT("/transactions/:id/paid", txnController.MarkAsPaid)
		r.PUT("/transactions/:id/extend", txnController.ExtendLoan)
	}

	// Dashboard Routes
	if dashController != nil {
		r.GET("/dashboard", dashController.GetSummary)
	}

	// Interest Routes
	if intController != nil {
		r.POST("/interest/calculate", intController.Calculate)
	}

	return r
}
