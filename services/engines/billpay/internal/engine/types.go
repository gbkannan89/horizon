package engine

type Biller struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	AccountNumber string `json:"account_number"`
	Amount      int64  `json:"amount"`
	DueDay      int    `json:"due_day"`
	AutoPay     bool   `json:"auto_pay"`
	Active      bool   `json:"active"`
	CreatedAt   string `json:"created_at"`
}

type Payment struct {
	ID        string `json:"id"`
	BillerID  string `json:"biller_id"`
	UserID    string `json:"user_id"`
	Amount    int64  `json:"amount"`
	Status    string `json:"status"`
	PaidAt    string `json:"paid_at,omitempty"`
}
