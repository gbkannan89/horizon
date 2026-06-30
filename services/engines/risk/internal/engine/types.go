package engine

// RiskLevel represents the severity of a risk score.
type RiskLevel string
const (
	RLCritical  RiskLevel = "Critical"
	RLElevated  RiskLevel = "Elevated"
	RLModerate  RiskLevel = "Moderate"
	RLLow       RiskLevel = "Low"
	RLMinimal   RiskLevel = "Minimal"
)

// RiskIndicator represents a single risk measurement.
type RiskIndicator struct {
	Name    string    `json:"name"`
	Score   int       `json:"score"`
	Level   RiskLevel `json:"level"`
	Warning int       `json:"warning_threshold"`
	Critical int      `json:"critical_threshold"`
}

// RiskCategoryScore represents a category-level aggregation.
type RiskCategoryScore struct {
	Category string    `json:"category"`
	Score    int       `json:"score"`
	Level    RiskLevel `json:"level"`
}

// StressTestResult represents a single stress scenario outcome.
type StressTestResult struct {
	Scenario     string `json:"scenario"`
	Impact       int64  `json:"impact"`
	RecoveryDays int    `json:"recovery_days"`
}

// RiskOutput is the complete risk assessment result.
type RiskOutput struct {
	CompositeScore int                 `json:"composite_score"`
	CompositeLevel RiskLevel           `json:"composite_level"`
	Indicators     []RiskIndicator     `json:"indicators"`
	CategoryScores []RiskCategoryScore `json:"category_scores"`
	StressTests    []StressTestResult  `json:"stress_tests"`
	Confidence     string              `json:"confidence"`
}
