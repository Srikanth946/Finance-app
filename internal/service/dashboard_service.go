package service

import (
	"finance_app/internal/models"
	"finance_app/internal/repository"
	"math"
	"strings"
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

	var activeUsers, paidUsers int
	var totalGiven, totalRecovered float64
	var previousMonthInterest float64
	var newUsersLastMonth int
	var recoveredLastMonth float64
	var projectedInterestLastMonth float64

	now := time.Now()
	firstDayOfCurrentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	firstDayOfLastMonth := firstDayOfCurrentMonth.AddDate(0, -1, 0)
	lastDayOfLastMonth := firstDayOfCurrentMonth.AddDate(0, 0, -1)

	// Get actual recovered amount from repository
	totalRecovered, _ = s.repo.GetTotalRecovered()

	for _, txn := range txns {
		totalGiven += txn.Amount

		// 1. Count New Users in last month
		if txn.CreatedAt.After(firstDayOfLastMonth) && txn.CreatedAt.Before(lastDayOfLastMonth) {
			newUsersLastMonth++
		}

		// FIXED: Case-insensitive status check
		if strings.ToUpper(txn.Status) == "ACTIVE" {
			activeUsers++

			summaryToday, err := s.interestSvc.CalculateInterest(txn)
			if err == nil && summaryToday.MonthsElapsed > 0 {
				// Approximation for the global report
				monthlyGrowth := (summaryToday.CurrentAmount - txn.Amount) / float64(summaryToday.MonthsElapsed)
				previousMonthInterest += monthlyGrowth
				projectedInterestLastMonth += monthlyGrowth
			}
		} else {
			paidUsers++
		}
	}

	// Recovered last month - utilizing the repository's month filter
	recoveredLastMonth, _ = s.repo.GetInterestForMonth(firstDayOfLastMonth.Year(), int(firstDayOfLastMonth.Month()))

	summary := &models.DashboardSummary{
		TotalUsers:            len(txns),
		ActiveUsers:           activeUsers,
		PaidUsers:             paidUsers,
		TotalAmountGiven:      totalGiven,
		TotalAmountRecovered:  totalRecovered,
		PreviousMonthInterest: math.Round(previousMonthInterest*100) / 100,
		NewUsersLastMonth:     newUsersLastMonth,
		RecoveredLastMonth:    math.Round(recoveredLastMonth*100) / 100,
		ProjectedInterest:     math.Round(projectedInterestLastMonth*100) / 100,
		LastUpdated:           now,
	}

	s.cache = summary
	s.cacheTime = now
	return summary, nil
}
