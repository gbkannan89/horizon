package api

import (
	"encoding/json"
	"net/http"

	"github.com/horizon/core/services/engines/billpay/internal/engine"
)

type Handler struct {
	store *engine.BillPayStore
}

func NewHandler(store *engine.BillPayStore) *Handler {
	return &Handler{store: store}
}

func userID(r *http.Request) string {
	if u := r.URL.Query().Get("user_id"); u != "" {
		return u
	}
	return "default"
}

func (h *Handler) AddBiller(w http.ResponseWriter, r *http.Request) {
	uid := userID(r)
	var req struct {
		Name          string `json:"name"`
		Category      string `json:"category"`
		AccountNumber string `json:"account_number"`
		Amount        int64  `json:"amount"`
		DueDay        int    `json:"due_day"`
		AutoPay       bool   `json:"auto_pay"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	b, err := h.store.AddBiller(uid, req.Name, req.Category, req.AccountNumber, req.Amount, req.DueDay, req.AutoPay)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": b})
}

func (h *Handler) ListBillers(w http.ResponseWriter, r *http.Request) {
	uid := userID(r)
	bill, err := h.store.ListBillers(uid)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if bill == nil { bill = []*engine.Biller{} }
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": bill})
}

func (h *Handler) GetBiller(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	b, err := h.store.GetBiller(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": b})
}

func (h *Handler) DeleteBiller(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.store.DeleteBiller(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handler) PayBill(w http.ResponseWriter, r *http.Request) {
	uid := userID(r)
	id := r.PathValue("id")
	billerID := id
	if billerID == "" {
		var req struct{ BillerID string `json:"biller_id"` }
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil && req.BillerID != "" {
			billerID = req.BillerID
		}
	}
	payment, err := h.store.PayBill(billerID, uid)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": payment})
}

func (h *Handler) PaymentHistory(w http.ResponseWriter, r *http.Request) {
	uid := userID(r)
	payments, err := h.store.ListPayments(uid)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if payments == nil { payments = []*engine.Payment{} }
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": payments})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
