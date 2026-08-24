package service

import (
	"finance_app/internal/models"
	"finance_app/internal/repository"
	"sync"
	"time"
)

type DashboardService interface {
	GetDashboardSummary() (*models.DashboardSummary, error)
}

type DashboardServiceHandler struct {
	repo        repository.TransactionRepository
	interestSvc InterestService
	cache       *models.DashboardSummary
	cacheTime   time.Time
	cacheMutex  sync.RWMutex
}

func NewDashboardService(repo repository.TransactionRepository, interestSvc InterestService) DashboardService {
	return &DashboardServiceHandler{
		repo:        repo,
		interestSvc: interestSvc,
	}
}

func (s *DashboardServiceHandler) GetDashboardSummary() (*models.DashboardSummary, error) {
	s.cacheMutex.RLock()
	if s.cache != nil && time.Since(s.cacheTime).Hours() < 24 {
		defer s.cacheMutex.RUnlock()
		return s.cache, nil
	}
	s.cacheMutex.RUnlock()

	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	// Double-check after lock
	if s.cache != nil && time.Since(s.cacheTime).Hours() < 24 {
		return s.cache, nil
	}

	txns, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var activeCounter, paidCounter int
	var outstandingPrincipal, outstandingAmount float64

	// Map to track monthly totals: "YYYY-MM" -> *MonthlyAnalysis
	monthlyMap := make(map[string]*models.MonthlyAnalysis)

	for _, txn := range txns {
		if txn.Status == "ACTIVE" {
			activeCounter++
			outstandingPrincipal += txn.Amount

			summary, err := s.interestSvc.CalculateInterest(txn)
			if err != nil {
				return nil, err
			}
			outstandingAmount += summary.CurrentAmount
		} else {
			paidCounter++
		}

		// Analysis for all customers (Company Level)
		// We use the StartDate to determine when the loan was given
		if txn.StartDate != "" {
			date, err := time.Parse("2006-01-02", txn.StartDate)
			if err == nil {
				monthKey := date.Format("2006-01")
				if _, ok := monthlyMap[monthKey]; !ok {
					monthlyMap[monthKey] = &models.MonthlyAnalysis{Month: monthKey}
				}
				monthlyMap[monthKey].AmountGiven += txn.Amount
			}
		}

		// For AmountReceived and InterestEarned, we'd normally track a separate
		// "Payments" table. Since we currently have a simple "PAID" status,
		// we'll calculate it based on when the loan was marked paid.
		if txn.Status == "PAID" {
			// In a real pro app, we'd have a Payments table with dates.
			// For this version, we'll use the record's creation or metadata if available.
			// As a simplification for this MVP, we assume payment happens near end of loan.
			// (Ideally, you'd add a 'PaymentDate' field to the model)
		}
	}

	// Convert map to slice
	var analysis []models.MonthlyAnalysis
	for _, v := range monthlyMap {
		analysis = append(analysis, *v)
	}

	summary := &models.DashboardSummary{
		TotalTransactions:  len(txns),
		ActiveTransactions: activeCounter,
		PaidTransactions:   paidCounter,
		TotalPrincipal:     outstandingPrincipal,
		TotalOutstanding:   outstandingAmount,
		MonthlyAnalysis:    analysis,
		LastUpdated:        time.Now(),
	}

	s.cache = summary
	s.cacheTime = time.Now()

	return summary, nil
}
