package main

import (
	"database/sql"
	"finance_app/internal/controller"
	"finance_app/internal/repository"
	"finance_app/internal/router"
	"finance_app/internal/service"
	"os"
	"time"

	"github.com/rs/zerolog"
	_ "modernc.org/sqlite"
)

func main() {
	// Configure Global Logger
	zerolog.TimeFieldFormat = time.RFC3339
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	log = log.With().
		Str("app", "finance_app").
		Str("env", "development").Logger()

	// Initialize Database
	db, err := sql.Open("sqlite", "finance.db")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open database")
	}
	defer db.Close()

	// Setup Repository
	txnRepo := repository.NewTransactionRepository(db)
	if err := txnRepo.InitializeDB(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}

	// Setup Service
	intService := service.NewInterestService(txnRepo)
	txnService := service.NewTransactionService(txnRepo, intService)
	dashService := service.NewDashboardService(txnRepo, intService)

	// Setup Controller
	txnController := controller.NewTransactionController(txnService)
	dashController := controller.NewDashboardController(dashService)
	intController := controller.NewInterestController(intService)

	// Setup Router
	r := router.SetupRouter(txnController, dashController, intController)

	log.Info().Msg("Server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
