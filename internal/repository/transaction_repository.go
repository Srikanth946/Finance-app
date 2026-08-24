package repository

import (
	"database/sql"
	"finance_app/internal/models"
	"fmt"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type TransactionRepository interface {
	Save(txn *models.Transaction) error
	UpdateStatus(mobileNumber string, status string) error
	UpdateAmount(mobileNumber string, amount float64) error
	Delete(mobileNumber string) error
	Exists(mobileNumber string) (bool, error)
	GetAll() ([]*models.Transaction, error)
	GetByMobile(mobileNumber string) (*models.Transaction, error)
	Search(surname, village, status *string) ([]*models.Transaction, error)
	InitializeDB() error
	SavePayment(payment *models.Payment) error
	GetTotalRecovered() (float64, error)
	GetInterestForMonth(year int, month int) (float64, error)
}

type sqliteTransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &sqliteTransactionRepository{db: db}
}

func (r *sqliteTransactionRepository) InitializeDB() error {
	query := `
    CREATE TABLE IF NOT EXISTS transactions(
        mobile_number TEXT PRIMARY KEY,
        surname TEXT NOT NULL,
        lastname TEXT NOT NULL,
        village TEXT NOT NULL,
        amount REAL NOT NULL,
        interest_rate REAL NOT NULL,
        start_date TEXT NOT NULL,
        end_date TEXT,
        compound_duration_months INTEGER NOT NULL,
        transaction_type TEXT NOT NULL,
        notes TEXT,
        status TEXT NOT NULL DEFAULT 'ACTIVE',
        created_at TEXT NOT NULL
    )`
	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	paymentQuery := `
    CREATE TABLE IF NOT EXISTS payments(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        mobile_number TEXT NOT NULL,
        amount_paid REAL NOT NULL,
        payment_date TEXT NOT NULL,
        payment_type TEXT NOT NULL
    )`
	_, err = r.db.Exec(paymentQuery)
	return err
}

func (r *sqliteTransactionRepository) Save(txn *models.Transaction) error {
	query := `
    INSERT INTO transactions(
        surname, lastname, village, mobile_number, amount, interest_rate, 
        start_date, end_date, compound_duration_months, 
        transaction_type, notes, status, created_at
    ) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)`

	_, err := r.db.Exec(query,
		strings.ToLower(txn.Surname),
		strings.ToLower(txn.LastName),
		txn.Village,
		txn.MobileNumber,
		txn.Amount,
		txn.InterestRate,
		txn.StartDate,
		txn.EndDate,
		txn.CompoundDurationMonths,
		txn.TransactionType,
		txn.Notes,
		strings.ToLower(txn.Status),
		txn.CreatedAt.Format(time.RFC3339),
	)
	return err
}

func (r *sqliteTransactionRepository) UpdateStatus(mobileNumber string, status string) error {
	query := `UPDATE transactions SET status = ? WHERE mobile_number = ?`
	_, err := r.db.Exec(query, strings.ToLower(status), mobileNumber)
	return err
}

func (r *sqliteTransactionRepository) UpdateAmount(mobileNumber string, amount float64) error {
	query := `UPDATE transactions SET amount = ? WHERE mobile_number = ?`
	_, err := r.db.Exec(query, amount, mobileNumber)
	return err
}

func (r *sqliteTransactionRepository) Delete(mobileNumber string) error {
	query := `DELETE FROM transactions WHERE mobile_number = ?`
	_, err := r.db.Exec(query, mobileNumber)
	return err
}

func (r *sqliteTransactionRepository) SavePayment(payment *models.Payment) error {
	query := `INSERT INTO payments(mobile_number, amount_paid, payment_date, payment_type) VALUES(?,?,?,?)`
	_, err := r.db.Exec(query, payment.MobileNumber, payment.AmountPaid, payment.PaymentDate.Format(time.RFC3339), payment.PaymentType)
	return err
}

func (r *sqliteTransactionRepository) GetTotalRecovered() (float64, error) {
	var total float64
	query := `SELECT SUM(amount_paid) FROM payments`
	err := r.db.QueryRow(query).Scan(&total)
	return total, err
}

func (r *sqliteTransactionRepository) GetInterestForMonth(year int, month int) (float64, error) {
	var total float64
	query := `SELECT SUM(amount_paid) FROM payments WHERE payment_type = 'INTEREST' AND strftime('%Y', payment_date) = ? AND strftime('%m', payment_date) = ?`
	monthStr := fmt.Sprintf("%02d", month)
	yearStr := fmt.Sprintf("%d", year)
	err := r.db.QueryRow(query, yearStr, monthStr).Scan(&total)
	return total, err
}

func (r *sqliteTransactionRepository) Exists(mobileNumber string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM transactions WHERE mobile_number = ? AND status = 'active')`
	err := r.db.QueryRow(query, mobileNumber).Scan(&exists)
	return exists, err
}

func (r *sqliteTransactionRepository) GetAll() ([]*models.Transaction, error) {
	query := `SELECT mobile_number, surname, lastname, village, amount, interest_rate, start_date, end_date, compound_duration_months, transaction_type, notes, status, created_at FROM transactions`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txns []*models.Transaction
	for rows.Next() {
		txn := &models.Transaction{}
		var createdAt string
		err := rows.Scan(&txn.MobileNumber, &txn.Surname, &txn.LastName, &txn.Village, &txn.Amount, &txn.InterestRate, &txn.StartDate, &txn.EndDate, &txn.CompoundDurationMonths, &txn.TransactionType, &txn.Notes, &txn.Status, &createdAt)
		if err != nil {
			return nil, err
		}
		txn.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		txns = append(txns, txn)
	}
	return txns, nil
}

func (r *sqliteTransactionRepository) GetByMobile(mobileNumber string) (*models.Transaction, error) {
	query := `SELECT mobile_number, surname, lastname, village, amount, interest_rate, start_date, end_date, compound_duration_months, transaction_type, notes, status, created_at FROM transactions WHERE mobile_number = ?`
	row := r.db.QueryRow(query, mobileNumber)

	txn := &models.Transaction{}
	var createdAt string
	err := row.Scan(&txn.MobileNumber, &txn.Surname, &txn.LastName, &txn.Village, &txn.Amount, &txn.InterestRate, &txn.StartDate, &txn.EndDate, &txn.CompoundDurationMonths, &txn.TransactionType, &txn.Notes, &txn.Status, &createdAt)
	if err != nil {
		return nil, err
	}
	txn.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return txn, nil
}

func (r *sqliteTransactionRepository) Search(surname, village, status *string) ([]*models.Transaction, error) {
	query := `SELECT mobile_number, surname, lastname, village, amount, interest_rate, start_date, end_date, compound_duration_months, transaction_type, notes, status, created_at FROM transactions WHERE 1=1`
	var args []interface{}

	if surname != nil {
		query += ` AND surname LIKE ?`
		args = append(args, "%"+*surname+"%")
	}
	if village != nil {
		query += ` AND village LIKE ?`
		args = append(args, "%"+*village+"%")
	}
	if status != nil {
		query += ` AND status = ?`
		args = append(args, *status)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txns []*models.Transaction
	for rows.Next() {
		txn := &models.Transaction{}
		var createdAt string
		err := rows.Scan(&txn.MobileNumber, &txn.Surname, &txn.LastName, &txn.Village, &txn.Amount, &txn.InterestRate, &txn.StartDate, &txn.EndDate, &txn.CompoundDurationMonths, &txn.TransactionType, &txn.Notes, &txn.Status, &createdAt)
		if err != nil {
			return nil, err
		}
		txn.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		txns = append(txns, txn)
	}
	return txns, nil
}
