package engine

// OptimizationStatus represents the lifecycle state.
type OptimizationStatus string
const (
	OSRequested  OptimizationStatus = "Requested"
	OSRunning    OptimizationStatus = "Running"
	OSCompleted  OptimizationStatus = "Completed"
	OSFailed     OptimizationStatus = "Failed"
	OSExpired    OptimizationStatus = "Expired"
	OSArchived   OptimizationStatus = "Archived"
)

// Objective represents a single optimization objective.
type Objective struct {
	Name  string  `json:"name"`
	Weight float64 `json:"weight"`
	Direction string `json:"direction"` // maximise, minimise
	Score float64 `json:"score"`
}

// Constraint represents a limitation on the solution.
type Constraint struct {
	Name    string `json:"name"`
	Type    string `json:"type"` // hard, soft
	Limit   float64 `json:"limit"`
}

// Solution represents a candidate optimal strategy.
type Solution struct {
	Rank           int                `json:"rank"`
	ObjectiveScores map[string]float64 `json:"objective_scores"`
	CompositeScore float64            `json:"composite_score"`
	Parameters     map[string]float64 `json:"parameters"`
	Description    string             `json:"description"`
}

// OptimizationOutput is the complete optimization result.
type OptimizationOutput struct {
	Objectives           []Objective  `json:"objectives"`
	Constraints          []Constraint `json:"constraints"`
	SelectedSolution     Solution     `json:"selected_solution"`
	Alternatives         []Solution   `json:"alternatives"`
	CandidatesEvaluated  int          `json:"candidates_evaluated"`
	CandidatesValid      int          `json:"candidates_valid"`
	Status               OptimizationStatus `json:"status"`
}
