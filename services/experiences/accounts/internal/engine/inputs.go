package engine

type Inputs struct {
	Accounts     []AccountInput `json:"accounts"`
	Account      *AccountInput  `json:"account,omitempty"`
	Transactions []TransInput   `json:"transactions,omitempty"`
	UserID       string         `json:"user_id"`
	SearchQuery  string         `json:"search_query,omitempty"`
}
