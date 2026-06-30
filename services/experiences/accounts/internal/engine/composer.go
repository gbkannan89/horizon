package engine

// Engine composes account experience view models from application service data.
type Engine struct{}

func NewEngine() *Engine { return &Engine{} }

// BuildAccountsList composes the paginated accounts dashboard.
func (e *Engine) BuildAccountsList(inputs Inputs, filter *FilterOpts, cursor string, limit int) *AccountsList {
	var cards []AccountCardData
	summary := &AccountsSummary{
		CountByType:   map[string]int{},
		CountByStatus: map[string]int{},
	}

	for _, a := range inputs.Accounts {
		if !e.matchesFilter(a, filter) {
			continue
		}
		cards = append(cards, AccountCardData{
			AccountID:        a.AccountID,
			AccountName:      a.AccountName,
			AccountType:      a.AccountType,
			Classification:   a.Classification,
			Status:           a.Status,
			Currency:         a.Currency,
			CurrentBalance:   a.CurrentBalance,
			AvailableBalance: a.AvailableBalance,
			InstitutionName:  a.InstitutionName,
			LastActivityDate: "",
			HasImport:        a.HasImport,
		})
		summary.TotalBalance += a.CurrentBalance
		summary.CountByType[a.AccountType]++
		summary.CountByStatus[a.Status]++
	}

	if limit <= 0 {
		limit = 25
	}
	hasMore := len(cards) > limit
	if len(cards) > limit {
		cards = cards[:limit]
	}

	return &AccountsList{
		Accounts: cards,
		Total:    len(inputs.Accounts),
		Cursor:   cursor,
		HasMore:  hasMore,
		Summary:  summary,
	}
}

// BuildAccountDetail composes the full account detail view.
func (e *Engine) BuildAccountDetail(a AccountInput, transactions []TransInput) *AccountDetailData {
	balances := e.buildBalances(a)

	var txn []TransPreview
	for _, t := range transactions {
		txn = append(txn, TransPreview{
			EventID:        t.EventID,
			EventType:      t.EventType,
			Amount:         t.Amount,
			Currency:       t.Currency,
			Description:    t.Description,
			Category:       t.Category,
			EventDate:      t.EventDate,
			RunningBalance: t.RunningBalance,
		})
	}
	if txn == nil {
		txn = []TransPreview{}
	}

	return &AccountDetailData{
		AccountID:        a.AccountID,
		AccountName:      a.AccountName,
		AccountType:      a.AccountType,
		Classification:   a.Classification,
		Status:           a.Status,
		Currency:         a.Currency,
		Balances:         balances,
		InstitutionName:  a.InstitutionName,
		ImportStatus:     a.ImportStatus,
		LastSyncTime:     a.LastSyncTime,
		Transactions:     txn,
		OpenedDate:       a.OpenedDate,
		ClosedDate:       a.ClosedDate,
		LiquidityProfile: a.LiquidityProfile,
		AccountHealth:    a.AccountHealth,
		CreditLimit:      a.CreditLimit,
		InterestRate:     a.InterestRate,
		Visibility:       a.Visibility,
		Tags:             a.Tags,
		Notes:            a.Notes,
		CreatedAt:        a.CreatedAt,
	}
}

func (e *Engine) buildBalances(a AccountInput) []BalanceEntry {
	return []BalanceEntry{
		{BTCurrent, a.CurrentBalance, "Current Balance", "Net sum of all posted transactions"},
		{BTAvailable, a.AvailableBalance, "Available Balance", "Current + uncleared − reserved"},
		{BTCleared, a.ClearedBalance, "Cleared Balance", "Confirmed and posted transactions only"},
		{BTUncleared, a.UnclearedBalance, "Uncleared Balance", "Pending transactions not yet confirmed"},
		{BTReserved, a.ReservedBalance, "Reserved Balance", "Funds committed but not yet allocated"},
		{BTAllocated, a.AllocatedBalance, "Allocated Balance", "Assigned to goals and allocations"},
		{BTSpendable, a.SpendableBalance, "Spendable Balance", "Available − allocated. Freely usable"},
	}
}

func (e *Engine) matchesFilter(a AccountInput, f *FilterOpts) bool {
	if f == nil {
		return true
	}
	if len(f.AccountTypes) > 0 && !inSlice(a.AccountType, f.AccountTypes) {
		return false
	}
	if len(f.Statuses) > 0 && !inSlice(a.Status, f.Statuses) {
		return false
	}
	if len(f.Currencies) > 0 && !inSlice(a.Currency, f.Currencies) {
		return false
	}
	return true
}

func inSlice(s string, list []string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
