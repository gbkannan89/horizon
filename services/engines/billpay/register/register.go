package register

import (
	"net/http"

	"github.com/horizon/core/services/engines/billpay/internal/api"
	"github.com/horizon/core/services/engines/billpay/internal/engine"
)

func RegisterRoutes(mux *http.ServeMux) {
	store := engine.NewBillPayStore()
	h := api.NewHandler(store)
	mux.HandleFunc("POST /api/v1/billpay/billers", h.AddBiller)
	mux.HandleFunc("GET /api/v1/billpay/billers", h.ListBillers)
	mux.HandleFunc("GET /api/v1/billpay/billers/{id}", h.GetBiller)
	mux.HandleFunc("DELETE /api/v1/billpay/billers/{id}", h.DeleteBiller)
	mux.HandleFunc("POST /api/v1/billpay/pay", h.PayBill)
	mux.HandleFunc("GET /api/v1/billpay/payments", h.PaymentHistory)
}
