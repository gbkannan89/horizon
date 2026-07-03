package http

import "net/http"

func RegisterBudgetHTTPRoutes(mux *http.ServeMux, h *BudgetHTTPHandler) {
	_ = h
}
