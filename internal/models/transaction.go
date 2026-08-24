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
	TotalUsers            int               `json:"total_users"`
	ActiveUsers           int               `json:"active_users"`
	PaidUsers             int               `json:"paid_users"`
	TotalAmountGiven      float64           `json:"total_amount_given"`
	TotalAmountRecovered  float64           `json:"total_amount_recovered"`
	PreviousMonthInterest float64           `json:"previous_month_interest"`
	NewUsersLastMonth     int               `json:"new_users_last_month"`
	RecoveredLastMonth    float64           `json:"recovered_last_month"`
	ProjectedInterest     float64           `json:"projected_interest_last_month"`
	MonthlyAnalysis       []MonthlyAnalysis `json:"monthly_analysis"`
	LastUpdated           time.Time         `json:"last_updated"`
}

type InterestSummary struct {
	OriginalAmount float64 `json:"original_amount"`
	CurrentAmount  float64 `json:"current_amount"`
	MonthsElapsed  int     `json:"months_elapsed"`
	TotalInterest  float64 `json:"total_interest"`
}

type Payment struct {
	ID           int       `json:"id"`
	MobileNumber string    `json:"mobile_number"`
	AmountPaid   float64   `json:"amount_paid"`
	PaymentDate  time.Time `json:"payment_date"`
	PaymentType  string    `json:"payment_type"` // PRINCIPAL, INTEREST, FULL_CLOSURE
}
