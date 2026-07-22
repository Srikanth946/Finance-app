package models

import "time"

type Transaction struct {
	ID                     int       `json:"id"`
	Surname                string    `json:"surname"`
	LastName               string    `json:"lastname"`
	Village                string    `json:"village"`
	Amount                 float64   `json:"amount"`
	InterestRate           float64   `json:"interest_rate"`
	StartDate              string    `json:"start_date"`
	EndDate                string    `json:"end_date,omitempty"`
	CompoundDurationMonths int       `json:"compound_duration_months"`
	TransactionType        string    `json:"transaction_type"`
	Notes                  string    `json:"notes,omitempty"`
	Status                 string    `json:"status"`
	CreatedAt              time.Time `json:"created_at"`
}

type InterestSummary struct {
	OriginalAmount     float64 `json:"original_amount"`
	CurrentAmount      float64 `json:"current_amount"`
	MonthsElapsed      int     `json:"months_elapsed"`
	TotalInterest      float64 `json:"total_interest"`
	CalculationEndDate string  `json:"calculation_end_date"`
}

type DashboardSummary struct {
	TotalTransactions  int     `json:"total_transactions"`
	ActiveTransactions int     `json:"active_transactions"`
	PaidTransactions   int     `json:"paid_transactions"`
	TotalPrincipal     float64 `json:"total_principal"`
	TotalOutstanding   float64 `json:"total_outstanding"`
}
