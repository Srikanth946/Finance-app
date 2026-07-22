package service

import (
	"finance_app/internal/models"
	"finance_app/internal/repository"
)

type DashboardService interface {
	GetDashboardSummary(interestSvc InterestService) (*models.DashboardSummary, error)
}

type DashboardServiceHandler struct {
	repo repository.TransactionRepository
}

func NewDashboardService(repo repository.TransactionRepository) DashboardService {
	return &DashboardServiceHandler{repo: repo}
}

func (s *DashboardServiceHandler) GetDashboardSummary(interestSvc InterestService) (*models.DashboardSummary, error) {
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

			summary, err := interestSvc.CalculateInterest(txn)
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
