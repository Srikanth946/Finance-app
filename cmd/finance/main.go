package main

import (
	"database/sql"
	"finance_app/internal/controller"
	"finance_app/internal/repository"
	"finance_app/internal/router"
	"finance_app/internal/service"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Configure Global Logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.With().
		Str("app", "finance_app").
		Str("env", "development").Logger()

	// Initialize Database
	db, err := sql.Open("sqlite3", "finance.db")
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
	txnService := service.NewTransactionService(txnRepo)

	// Setup Controller
	txnController := controller.NewTransactionController(txnService)

	// Setup Router
	r := router.SetupRouter(txnController)

	log.Info().Msg("Server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
