package engine

// SimStatus represents the simulation lifecycle state.
type SimStatus string
const (
	SSDraft     SimStatus = "Draft"
	SSRunning   SimStatus = "Running"
	SSCompleted SimStatus = "Completed"
	SSFailed    SimStatus = "Failed"
	SSExpired   SimStatus = "Expired"
	SSArchived  SimStatus = "Archived"
	SSHistorical SimStatus = "Historical"
)

// SimType represents the type of simulation.
type SimType string
const (
	STIncreaseSIP        SimType = "IncreaseSIP"
	STDecreaseSIP        SimType = "DecreaseSIP"
	STEarlyLoanClosure   SimType = "EarlyLoanClosure"
	STExtraEMI           SimType = "ExtraEMI"
	STSalaryIncrease     SimType = "SalaryIncrease"
	STSalaryLoss         SimType = "SalaryLoss"
	STJobChange          SimType = "JobChange"
	STMarriage           SimType = "Marriage"
	STChild              SimType = "Child"
	STHousePurchase      SimType = "HousePurchase"
	STVehiclePurchase    SimType = "VehiclePurchase"
	STInvestmentLumpSum  SimType = "InvestmentLumpSum"
	STExpenseReduction   SimType = "ExpenseReduction"
	STGoalDelay          SimType = "GoalDelay"
	STGoalPriorityChange SimType = "GoalPriorityChange"
	STEmergency          SimType = "Emergency"
	STMarketCrash        SimType = "MarketCrash"
	STMarketBoom         SimType = "MarketBoom"
	STInflationChange    SimType = "InflationChange"
	STRetirementAgeChange SimType = "RetirementAgeChange"
)

// DeltaEntry represents a single metric delta.
type DeltaEntry struct {
	Metric        string `json:"metric"`
	BaselineValue string `json:"baseline_value"`
	ScenarioValue string `json:"scenario_value"`
	Delta         string `json:"delta"`
	Direction     string `json:"direction"` // Improvement, Regression, NoChange
}

// SimulationOutput is the result of a simulation.
type SimulationOutput struct {
	ScenarioID      string       `json:"scenario_id"`
	SimType         SimType      `json:"sim_type"`
	Status          SimStatus    `json:"status"`
	BaselineVersion string       `json:"baseline_version"`
	DeltaSummary    []DeltaEntry `json:"delta_summary"`
	Parameters      map[string]interface{} `json:"parameters"`
	Confidence      string       `json:"confidence"`
}
