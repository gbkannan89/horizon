package domain

type Frequency string

const (
	FreqDaily      Frequency = "Daily"
	FreqWeekly     Frequency = "Weekly"
	FreqBiWeekly   Frequency = "BiWeekly"
	FreqMonthly    Frequency = "Monthly"
	FreqQuarterly  Frequency = "Quarterly"
	FreqSemiAnnual Frequency = "SemiAnnual"
	FreqAnnual     Frequency = "Annual"
	FreqCustom     Frequency = "Custom"
)

type RecurringStatus string

const (
	StatusActive    RecurringStatus = "Active"
	StatusPaused    RecurringStatus = "Paused"
	StatusCompleted RecurringStatus = "Completed"
	StatusCancelled RecurringStatus = "Cancelled"
	StatusArchived  RecurringStatus = "Archived"
)

type OccurrenceResult string

const (
	ResultGenerated OccurrenceResult = "Generated"
	ResultSkipped   OccurrenceResult = "Skipped"
	ResultFailed    OccurrenceResult = "Failed"
)
