package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const BaseURL = "http://localhost:8080"

type DashboardSummary struct {
	PreviousMonthInterest float64 `json:"previous_month_interest"`
	TotalAmount           float64 `json:"total_amount_given"`
}

func GetDashboardSummary(ctx context.Context) (string, error) {
	resp, err := http.Get(BaseURL + "/dashboard")
	if err != nil {
		return "", fmt.Errorf("failed to call dashboard: %v", err)
	}
	defer resp.Body.Close()

	var summary DashboardSummary
	if err := json.NewDecoder(resp.Body).Decode(&summary); err != nil {
		return "", fmt.Errorf("failed to decode summary: %v", err)
	}

	return fmt.Sprintf("Dashboard Summary: Previous Month Interest: %.2f, Total Amount: %.2f",
		summary.PreviousMonthInterest, summary.TotalAmount), nil
}

func GetTransactions(ctx context.Context) (string, error) {
	resp, err := http.Get(BaseURL + "/transactions")
	if err != nil {
		return "", fmt.Errorf("failed to fetch transactions: %v", err)
	}
	defer resp.Body.Close()

	var txns []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&txns); err != nil {
		return "", fmt.Errorf("failed to decode transactions: %v", err)
	}

	return fmt.Sprintf("Found %d transactions", len(txns)), nil
}
