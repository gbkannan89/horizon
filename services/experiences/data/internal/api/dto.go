package api

import "time"

type ExportData struct {
	Version       string                     `json:"version"`
	ExportedAt    string                     `json:"exported_at"`
	UserID        string                     `json:"user_id"`
	Accounts      []map[string]interface{}   `json:"accounts"`
	Transactions  []map[string]interface{}   `json:"transactions"`
	Goals         []map[string]interface{}   `json:"goals"`
	Allocations   []map[string]interface{}   `json:"allocations"`
	Assets        []map[string]interface{}   `json:"assets"`
	Liabilities   []map[string]interface{}   `json:"liabilities"`
	Portfolios    []map[string]interface{}   `json:"portfolios"`
	Institutions  []map[string]interface{}   `json:"institutions"`
	Households    []map[string]interface{}   `json:"households"`
	HealthScores  []map[string]interface{}   `json:"health_scores"`
	RiskAssessments []map[string]interface{} `json:"risk_assessments"`
	Recommendations []map[string]interface{} `json:"recommendations"`
	Optimizations []map[string]interface{}   `json:"optimizations"`
	Projections   []map[string]interface{}   `json:"projections"`
	Simulations   []map[string]interface{}   `json:"simulations"`
}

func NewExportData(userID string) ExportData {
	return ExportData{
		Version:    "1.0",
		ExportedAt: time.Now().UTC().Format(time.RFC3339),
		UserID:     userID,
		Accounts:      []map[string]interface{}{},
		Transactions:  []map[string]interface{}{},
		Goals:         []map[string]interface{}{},
		Allocations:   []map[string]interface{}{},
		Assets:        []map[string]interface{}{},
		Liabilities:   []map[string]interface{}{},
		Portfolios:    []map[string]interface{}{},
		Institutions:  []map[string]interface{}{},
		Households:    []map[string]interface{}{},
		HealthScores:  []map[string]interface{}{},
		RiskAssessments: []map[string]interface{}{},
		Recommendations: []map[string]interface{}{},
		Optimizations: []map[string]interface{}{},
		Projections:   []map[string]interface{}{},
		Simulations:   []map[string]interface{}{},
	}
}
