package service

import (
	"finance_app/internal/models"
	"finance_app/internal/repository"
	"math"
	"time"
)

type InterestService interface {
	CalculateInterest(txn *models.Transaction) (*models.InterestSummary, error)
	CalculateGeneralInterest(amount float64, rate float64, months int, compoundDuration int) (*models.InterestSummary, error)
}

type InterestServiceHandler struct {
	repo repository.TransactionRepository
}

func NewInterestService(repo repository.TransactionRepository) *InterestServiceHandler {
	return &InterestServiceHandler{repo: repo}
}

func (s *InterestServiceHandler) CalculateInterest(txn *models.Transaction) (*models.InterestSummary, error) {
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

	// Calculate exact days elapsed for pro-rata precision
	daysElapsed := int(calculationEndDate.Sub(startDate).Hours() / 24)
	if daysElapsed < 0 {
		daysElapsed = 0
	}
	monthsElapsed := daysElapsed / 30 // Approximate months for the summary report

	principal := txn.Amount
	dailyRate := txn.InterestRate / 100 / 365
	remainingDays := daysElapsed
	totalInterest := 0.0

	// Compound based on CompoundDurationMonths (approx 30 days per month)
	compoundDays := txn.CompoundDurationMonths * 30

	for remainingDays > 0 {
		daysToCalculate := remainingDays
		if daysToCalculate > compoundDays {
			daysToCalculate = compoundDays
		}

		interest := principal * dailyRate * float64(daysToCalculate)
		totalInterest += interest
		principal += interest
		remainingDays -= daysToCalculate
	}

	return &models.InterestSummary{
		OriginalAmount: txn.Amount,
		CurrentAmount:  math.Round(principal*100) / 100,
		MonthsElapsed:  monthsElapsed,
		TotalInterest:  math.Round(totalInterest*100) / 100,
	}, nil
}

func (s *InterestServiceHandler) CalculateGeneralInterest(amount float64, rate float64, months int, compoundDuration int) (*models.InterestSummary, error) {
	principal := amount
	monthlyRate := rate / 100 / 12
	remainingMonths := months
	totalInterest := 0.0

	for remainingMonths > 0 {
		monthsToCalculate := remainingMonths
		if monthsToCalculate > compoundDuration {
			monthsToCalculate = compoundDuration
		}

		interest := principal * monthlyRate * float64(monthsToCalculate)
		totalInterest += interest
		principal += interest
		remainingMonths -= monthsToCalculate
	}

	return &models.InterestSummary{
		OriginalAmount: amount,
		CurrentAmount:  math.Round(principal*100) / 100,
		MonthsElapsed:  months,
		TotalInterest:  math.Round(totalInterest*100) / 100,
	}, nil
}
