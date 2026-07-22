package service

import (
	"finance_app/internal/models"
	"finance_app/internal/repository"
	"time"
)

type TransactionService interface {
	CreateTransaction(txn *models.Transaction) error
	GetAllTransactions() ([]*models.Transaction, error)
	GetTransactionByID(id int) (*models.Transaction, error)
	MarkAsPaid(id int) error
	FilterTransactions(surname, village, status *string) ([]*models.Transaction, error)
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
