package engine

// ScoreGrade represents the overall health score grade.
type ScoreGrade string
const (SGExcellent ScoreGrade = "Excellent"; SGGood ScoreGrade = "Good"; SGFair ScoreGrade = "Fair"; SGNeedsWork ScoreGrade = "NeedsWork"; SGCritical ScoreGrade = "Critical")

// Trend represents the score trend direction.
type Trend string
const (TImproving Trend = "Improving"; TStable Trend = "Stable"; TDeclining Trend = "Declining"; TVolatile Trend = "Volatile")

// DimensionScore represents a single health dimension score.
type DimensionScore struct {
	Name         string  `json:"name"`
	Score        float64 `json:"score"`
	Weight       float64 `json:"weight"`
	Contribution float64 `json:"contribution"`
}

// HealthScoreOutput is the result of a health score computation.
type HealthScoreOutput struct {
	ScoreID       string           `json:"score_id"`
	OverallScore  int              `json:"overall_score"`
	ScoreGrade    ScoreGrade       `json:"score_grade"`
	Trend         Trend            `json:"trend"`
	Dimensions    []DimensionScore `json:"dimensions"`
	PreviousScore int              `json:"previous_score"`
	Confidence    string           `json:"confidence"`
}
