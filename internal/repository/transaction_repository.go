package repository

import (
	"database/sql"
	"finance_app/internal/models"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type TransactionRepository interface {
	Save(txn *models.Transaction) error
	UpdateStatus(id int, status string) error
	GetAll() ([]*models.Transaction, error)
	GetByID(id int) (*models.Transaction, error)
	Search(surname, village, status *string) ([]*models.Transaction, error)
	InitializeDB() error
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
        id INTEGER PRIMARY KEY AUTOINCREMENT,
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
	return err
}

func (r *sqliteTransactionRepository) Save(txn *models.Transaction) error {
	query := `
    INSERT INTO transactions(
        surname, lastname, village, amount, interest_rate, 
        start_date, end_date, compound_duration_months, 
        transaction_type, notes, status, created_at
    ) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)`

	_, err := r.db.Exec(query,
		strings.ToLower(txn.Surname),
		strings.ToLower(txn.LastName),
		txn.Village,
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

func (r *sqliteTransactionRepository) UpdateStatus(id int, status string) error {
	query := `UPDATE transactions SET status = ? WHERE id = ?`
	_, err := r.db.Exec(query, strings.ToLower(status), id)
	return err
}

func (r *sqliteTransactionRepository) GetAll() ([]*models.Transaction, error) {
	query := `SELECT * FROM transactions`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txns []*models.Transaction
	for rows.Next() {
		txn := &models.Transaction{}
		var createdAtStr string
		err := rows.Scan(&txn.ID, &txn.Surname, &txn.LastName, &txn.Village, &txn.Amount,
			&txn.InterestRate, &txn.StartDate, &txn.EndDate, &txn.CompoundDurationMonths,
			&txn.TransactionType, &txn.Notes, &txn.Status, &createdAtStr)
		if err != nil {
			return nil, err
		}
		txn.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		txns = append(txns, txn)
	}
	return txns, nil
}

func (r *sqliteTransactionRepository) GetByID(id int) (*models.Transaction, error) {
	query := `SELECT * FROM transactions WHERE id = ?`
	row := r.db.QueryRow(query, id)

	txn := &models.Transaction{}
	var createdAtStr string
	err := row.Scan(&txn.ID, &txn.Surname, &txn.LastName, &txn.Village, &txn.Amount,
		&txn.InterestRate, &txn.StartDate, &txn.EndDate, &txn.CompoundDurationMonths,
		&txn.TransactionType, &txn.Notes, &txn.Status, &createdAtStr)
	if err != nil {
		return nil, err
	}
	txn.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
	return txn, nil
}

func (r *sqliteTransactionRepository) Search(surname, village, status *string) ([]*models.Transaction, error) {
	query := `SELECT * FROM transactions WHERE 1=1`
	var args []interface{}

	if surname != nil {
		query += ` AND surname = ?`
		args = append(args, strings.ToLower(*surname))
	}
	if village != nil {
		query += ` AND village = ?`
		args = append(args, *village)
	}
	if status != nil {
		query += ` AND status = ?`
		args = append(args, strings.ToLower(*status))
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txns []*models.Transaction
	for rows.Next() {
		txn := &models.Transaction{}
		var createdAtStr string
		err := rows.Scan(&txn.ID, &txn.Surname, &txn.LastName, &txn.Village, &txn.Amount,
			&txn.InterestRate, &txn.StartDate, &txn.EndDate, &txn.CompoundDurationMonths,
			&txn.TransactionType, &txn.Notes, &txn.Status, &createdAtStr)
		if err != nil {
			return nil, err
		}
		txn.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		txns = append(txns, txn)
	}
	return txns, nil
}
