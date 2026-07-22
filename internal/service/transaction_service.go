package service

import (
	"finance_app/internal/models"
	"finance_app/internal/repository"
	"math"
	"time"
)

type TransactionService interface {
	CreateTransaction(txn *models.Transaction) error
	GetAllTransactions() ([]*models.Transaction, error)
	GetTransactionByID(id int) (*models.Transaction, error)
	MarkAsPaid(id int) error
	FilterTransactions(surname, village, status *string) ([]*models.Transaction, error)
	CalculateInterest(txn *models.Transaction) (*models.InterestSummary, error)
	GetDashboardSummary() (*models.DashboardSummary, error)
}

type transactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{repo: repo}
}

func (s *transactionService) CreateTransaction(txn *models.Transaction) error {
	txn.CreatedAt = time.Now()
	return s.repo.Save(txn)
}

func (s *transactionService) GetAllTransactions() ([]*models.Transaction, error) {
	return s.repo.GetAll()
}

func (s *transactionService) GetTransactionByID(id int) (*models.Transaction, error) {
	return s.repo.GetByID(id)
}

func (s *transactionService) MarkAsPaid(id int) error {
	return s.repo.UpdateStatus(id, "PAID")
}

func (s *transactionService) FilterTransactions(surname, village, status *string) ([]*models.Transaction, error) {
	return s.repo.Search(surname, village, status)
}

func (s *transactionService) CalculateInterest(txn *models.Transaction) (*models.InterestSummary, error) {
	var calculationEndDate time.Time
	now := time.Now()

	if txn.EndDate != "" {
		t, err := time.Parse("2006-01-02", txn.EndDate)
		if err != nil {
			return nil, err
		}
		calculationEndDate = t
	} else {
		calculationEndDate = now
	}

	startDate, err := time.Parse("2006-01-02", txn.StartDate)
	if err != nil {
		return nil, err
	}

	// Calculate months elapsed
	years := calculationEndDate.Year() - startDate.Year()
	months := int(calculationEndDate.Month()) - int(startDate.Month())
	monthsElapsed := years*12 + months

	if calculationEndDate.Day() < startDate.Day() {
		monthsElapsed--
	}

	principal := txn.Amount
	monthlyRate := txn.InterestRate / 100 / 12
	remainingMonths := monthsElapsed
	totalInterest := 0.0

	for remainingMonths > 0 {
		monthsToCalculate := remainingMonths
		if monthsToCalculate > txn.CompoundDurationMonths {
			monthsToCalculate = txn.CompoundDurationMonths
		}

		interest := principal * monthlyRate * float64(monthsToCalculate)
		totalInterest += interest
		principal += interest
		remainingMonths -= monthsToCalculate
	}

	return &models.InterestSummary{
		OriginalAmount:     txn.Amount,
		CurrentAmount:      math.Round(principal*100) / 100,
		MonthsElapsed:      monthsElapsed,
		TotalInterest:      math.Round(totalInterest*100) / 100,
		CalculationEndDate: calculationEndDate.Format("2006-01-02"),
	}, nil
}

func (s *transactionService) GetDashboardSummary() (*models.DashboardSummary, error) {
	txns, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var activeCounter, paidCounter int
	var outstandingPrincipal, outstandingAmount float64

	for _, txn := range txns {
		if txn.Status == "ACTIVE" {
			activeCounter++
			outstandingPrincipal += txn.Amount

			summary, err := s.CalculateInterest(txn)
			if err != nil {
				return nil, err
			}
			outstandingAmount += summary.CurrentAmount
		} else {
			paidCounter++
		}
	}

	return &models.DashboardSummary{
		TotalTransactions:  len(txns),
		ActiveTransactions: activeCounter,
		PaidTransactions:   paidCounter,
		TotalPrincipal:     outstandingPrincipal,
		TotalOutstanding:   outstandingAmount,
	}, nil
}
