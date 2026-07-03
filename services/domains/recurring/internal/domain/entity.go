package domain

import (
	"time"
)

func ReconstructFromDB(
	id, userID, householdID, name, description string,
	amount int64, currency string,
	freq string, interval int,
	startDate time.Time, endDate, nextOccurrence, lastOccurrence *time.Time,
	status string,
	eventType, category, source, destination string,
	skipHolidays, skipWeekends bool,
	tags []string, metadata map[string]string,
	createdAt, updatedAt time.Time,
) *RecurringTransaction {
	if tags == nil { tags = []string{} }
	if metadata == nil { metadata = map[string]string{} }
	return &RecurringTransaction{
		id: id, userID: userID, householdID: householdID,
		name: name, description: description,
		amount: amount, currency: currency,
		frequency: Frequency(freq), interval: interval,
		startDate: startDate, endDate: endDate,
		nextOccurrence: nextOccurrence, lastOccurrence: lastOccurrence,
		status: RecurringStatus(status),
		eventTemplate: EventTemplate{Type: eventType, Category: category, Source: source, Destination: destination},
		skipHolidays: skipHolidays, skipWeekends: skipWeekends,
		tags: tags, metadata: metadata,
		createdAt: createdAt, updatedAt: updatedAt,
	}
}

type EventTemplate struct {
	Type        string `json:"type"`
	Category    string `json:"category"`
	Source      string `json:"source"`
	Destination string `json:"destination,omitempty"`
}

type RecurringTransaction struct {
	id             string
	userID         string
	householdID    string
	name           string
	description    string
	amount         int64
	currency       string
	frequency      Frequency
	interval       int
	startDate      time.Time
	endDate        *time.Time
	nextOccurrence *time.Time
	lastOccurrence *time.Time
	status         RecurringStatus
	eventTemplate  EventTemplate
	skipHolidays   bool
	skipWeekends   bool
	tags           []string
	metadata       map[string]string
	createdAt      time.Time
	updatedAt      time.Time
}

func (r *RecurringTransaction) ID() string { return r.id }
func (r *RecurringTransaction) UserID() string { return r.userID }
func (r *RecurringTransaction) HouseholdID() string { return r.householdID }
func (r *RecurringTransaction) Name() string { return r.name }
func (r *RecurringTransaction) Description() string { return r.description }
func (r *RecurringTransaction) Amount() int64 { return r.amount }
func (r *RecurringTransaction) Currency() string { return r.currency }
func (r *RecurringTransaction) Frequency() Frequency { return r.frequency }
func (r *RecurringTransaction) Interval() int { return r.interval }
func (r *RecurringTransaction) StartDate() time.Time { return r.startDate }
func (r *RecurringTransaction) EndDate() *time.Time { return r.endDate }
func (r *RecurringTransaction) NextOccurrence() *time.Time { return r.nextOccurrence }
func (r *RecurringTransaction) LastOccurrence() *time.Time { return r.lastOccurrence }
func (r *RecurringTransaction) Status() RecurringStatus { return r.status }
func (r *RecurringTransaction) EventTemplate() EventTemplate { return r.eventTemplate }
func (r *RecurringTransaction) SkipHolidays() bool { return r.skipHolidays }
func (r *RecurringTransaction) SkipWeekends() bool { return r.skipWeekends }
func (r *RecurringTransaction) Tags() []string { return r.tags }
func (r *RecurringTransaction) Metadata() map[string]string { return r.metadata }
func (r *RecurringTransaction) CreatedAt() time.Time { return r.createdAt }
func (r *RecurringTransaction) UpdatedAt() time.Time { return r.updatedAt }
