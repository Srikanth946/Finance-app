package service

import (
	"finance_app/internal/models"
	"finance_app/internal/repository"
	"fmt"
	"time"
)

type TransactionService interface {
	CreateTransaction(txn *models.Transaction) error
	GetAllTransactions() ([]*models.Transaction, error)
	GetTransactionByMobile(mobileNumber string) (*models.Transaction, error)
	GetLoanWithInterest(mobileNumber string) (*models.Transaction, *models.InterestSummary, error)
	MarkAsPaid(mobileNumber string) error
	FilterTransactions(surname, village, status *string) ([]*models.Transaction, error)
	ExtendLoan(mobileNumber string, newRate float64, newCompoundDuration int) error
	GetSortedTransactions(sortBy string) ([]*models.Transaction, error)
}

type transactionService struct {
	repo        repository.TransactionRepository
	interestSvc InterestService
}

func NewTransactionService(repo repository.TransactionRepository, interestSvc InterestService) TransactionService {
	return &transactionService{
		repo:        repo,
		interestSvc: interestSvc,
	}
}

func (s *transactionService) CreateTransaction(txn *models.Transaction) error {
	// If start date is provided, use it as the creation date
	if txn.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", txn.StartDate)
		if err == nil {
			txn.CreatedAt = startDate
		} else {
			txn.CreatedAt = time.Now()
		}
	} else {
		txn.CreatedAt = time.Now()
	}

	exists, err := s.repo.Exists(txn.MobileNumber)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("a transaction for this mobile number already exists and is active")
	}

	return s.repo.Save(txn)
}

func (s *transactionService) GetAllTransactions() ([]*models.Transaction, error) {
	return s.repo.GetAll()
}

func (s *transactionService) GetTransactionByMobile(mobileNumber string) (*models.Transaction, error) {
	return s.repo.GetByMobile(mobileNumber)
}

func (s *transactionService) GetLoanWithInterest(mobileNumber string) (*models.Transaction, *models.InterestSummary, error) {
	txn, err := s.repo.GetByMobile(mobileNumber)
	if err != nil {
		return nil, nil, err
	}

	summary, err := s.interestSvc.CalculateInterest(txn)
	if err != nil {
		return txn, nil, err
	}

	return txn, summary, nil
}

func (s *transactionService) MarkAsPaid(mobileNumber string) error {
	return s.repo.Delete(mobileNumber)
}

func (s *transactionService) FilterTransactions(surname, village, status *string) ([]*models.Transaction, error) {
	return s.repo.Search(surname, village, status)
}

func (s *transactionService) ExtendLoan(mobileNumber string, newRate float64, newCompoundDuration int) error {
	txn, err := s.repo.GetByMobile(mobileNumber)
	if err != nil {
		return err
	}

	if txn.Status == "PAID" {
		return fmt.Errorf("transaction is already paid")
	}

	summary, err := s.interestSvc.CalculateInterest(txn)
	if err != nil {
		return err
	}

	// 1. Delete current loan to free up the mobile number unique constraint
	if err := s.repo.Delete(mobileNumber); err != nil {
		return fmt.Errorf("failed to remove old loan before extension: %v", err)
	}

	// 2. Create new loan with current balance as principal and NEW terms
	newTxn := *txn
	newTxn.Amount = summary.CurrentAmount
	newTxn.InterestRate = newRate
	newTxn.CompoundDurationMonths = newCompoundDuration
	newTxn.StartDate = time.Now().Format("2006-01-02")
	newTxn.Status = "ACTIVE"
	newTxn.Notes = fmt.Sprintf("Extended from loan with mobile %s. New Rate: %.2f%%, Duration: %d", mobileNumber, newRate, newCompoundDuration)

	return s.repo.Save(&newTxn)
}

func (s *transactionService) GetSortedTransactions(sortBy string) ([]*models.Transaction, error) {
	txns, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	// ... implementation below ...
	return txns, nil
}
