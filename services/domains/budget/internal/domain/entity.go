package domain

import "time"

type Budget struct {
	id              string
	userID          string
	householdID     string
	name            string
	period          BudgetPeriod
	startDate       time.Time
	endDate         time.Time
	status          BudgetStatus
	totalBudgeted   int64
	totalSpent      int64
	totalRemaining  int64
	currency        string
	categories      []BudgetCategory
	tags            []string
	metadata        map[string]string
	createdAt       time.Time
	updatedAt       time.Time
}

func (b *Budget) ID() string                    { return b.id }
func (b *Budget) UserID() string                { return b.userID }
func (b *Budget) HouseholdID() string           { return b.householdID }
func (b *Budget) Name() string                  { return b.name }
func (b *Budget) Period() BudgetPeriod          { return b.period }
func (b *Budget) StartDate() time.Time          { return b.startDate }
func (b *Budget) EndDate() time.Time            { return b.endDate }
func (b *Budget) Status() BudgetStatus          { return b.status }
func (b *Budget) TotalBudgeted() int64          { return b.totalBudgeted }
func (b *Budget) TotalSpent() int64             { return b.totalSpent }
func (b *Budget) TotalRemaining() int64         { return b.totalRemaining }
func (b *Budget) Currency() string              { return b.currency }
func (b *Budget) Categories() []BudgetCategory  { return b.categories }
func (b *Budget) Tags() []string                { return b.tags }
func (b *Budget) Metadata() map[string]string   { return b.metadata }
func (b *Budget) CreatedAt() time.Time          { return b.createdAt }
func (b *Budget) UpdatedAt() time.Time          { return b.updatedAt }

func (b *Budget) SetStatus(s BudgetStatus) {
	b.status = s
	b.updatedAt = time.Now().UTC()
}

func ReconstructFromDB(
	id, userID, householdID, name string,
	period BudgetPeriod,
	startDate, endDate time.Time,
	status BudgetStatus,
	totalBudgeted, totalSpent, totalRemaining int64,
	currency string,
	categories []BudgetCategory,
	tags []string,
	createdAt, updatedAt time.Time,
) *Budget {
	if tags == nil { tags = []string{} }
	if categories == nil { categories = []BudgetCategory{} }
	return &Budget{
		id: id, userID: userID, householdID: householdID,
		name: name, period: period, startDate: startDate, endDate: endDate,
		status: status, totalBudgeted: totalBudgeted, totalSpent: totalSpent,
		totalRemaining: totalRemaining, currency: currency,
		categories: categories, tags: tags, metadata: map[string]string{},
		createdAt: createdAt, updatedAt: updatedAt,
	}
}
