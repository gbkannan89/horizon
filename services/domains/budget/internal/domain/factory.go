package domain

import (
	"time"
)

type BudgetFactory struct{}

func NewBudgetFactory() *BudgetFactory { return &BudgetFactory{} }

func (f *BudgetFactory) Create(
	id, userID, householdID, name string,
	period BudgetPeriod,
	startDate, endDate time.Time,
	currency string,
	categories []BudgetCategory,
	tags []string,
) (*Budget, error) {

	if err := ValidateBudgetName(name); err != nil { return nil, err }
	if err := ValidateBudgetPeriod(period); err != nil { return nil, err }
	if err := ValidateDateRange(startDate, endDate); err != nil { return nil, err }
	if err := ValidateCurrency(currency); err != nil { return nil, err }
	for _, c := range categories {
		if err := ValidateCategoryName(c.Category); err != nil { return nil, err }
		if err := ValidateBudgetAmount(c.BudgetedAmount); err != nil { return nil, err }
	}

	if tags == nil { tags = []string{} }
	if categories == nil { categories = []BudgetCategory{} }

	var totalBudgeted int64
	for _, c := range categories {
		totalBudgeted += c.BudgetedAmount
	}

	now := time.Now().UTC()
	return &Budget{
		id: id, userID: userID, householdID: householdID,
		name: name, period: period, status: BStatusDraft,
		startDate: startDate, endDate: endDate,
		totalBudgeted: totalBudgeted, totalSpent: 0, totalRemaining: totalBudgeted,
		currency: currency, categories: categories,
		tags: tags, metadata: map[string]string{},
		createdAt: now, updatedAt: now,
	}, nil
}
