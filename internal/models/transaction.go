package models

import "time"

type Transaction struct {
	MobileNumber           string    `json:"mobile_number"`
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
	Status                 string    `json:"status"` // ACTIVE, PAID
	CreatedAt              time.Time `json:"created_at"`
}

type MonthlyAnalysis struct {
	Month          string  `json:"month"`
	AmountGiven    float64 `json:"amount_given"`
	AmountReceived float64 `json:"amount_received"`
	InterestEarned float64 `json:"interest_earned"`
}

type DashboardSummary struct {
	TotalTransactions  int               `json:"total_transactions"`
	ActiveTransactions int               `json:"active_transactions"`
	PaidTransactions   int               `json:"paid_transactions"`
	TotalPrincipal     float64           `json:"total_principal"`
	TotalOutstanding   float64           `json:"total_outstanding"`
	MonthlyAnalysis    []MonthlyAnalysis `json:"monthly_analysis"`
	LastUpdated        time.Time         `json:"last_updated"`
}

type InterestSummary struct {
	OriginalAmount float64 `json:"original_amount"`
	CurrentAmount  float64 `json:"current_amount"`
	MonthsElapsed  int     `json:"months_elapsed"`
	TotalInterest  float64 `json:"total_interest"`
}
