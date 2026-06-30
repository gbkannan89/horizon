package engine

// InputReader reads domain state for the engine.
type InputReader interface {
	ReadProjectionInputs() (*Inputs, error)
}

// Config holds engine configuration.
type Config struct {
	MaxProjectionHorizon    int `json:"max_projection_horizon"`   // years
	SnapshotIntervalMonths  int `json:"snapshot_interval_months"` // months
	ConfidenceLookbackYears int `json:"confidence_lookback_years"`
}

// DefaultConfig returns default engine configuration.
func DefaultConfig() Config {
	return Config{
		MaxProjectionHorizon: 30,
		SnapshotIntervalMonths: 12,
		ConfidenceLookbackYears: 3,
	}
}
