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
