package http

import "github.com/horizon/core/services/domains/budget/internal/application"

type BudgetHTTPHandler struct {
	svc *application.BudgetService
}

func NewBudgetHTTPHandler(svc *application.BudgetService) *BudgetHTTPHandler {
	return &BudgetHTTPHandler{svc: svc}
}
