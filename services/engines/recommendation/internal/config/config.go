package engine

// Config holds recommendation engine configuration.
type Config struct {
	MaxRecommendationsPerCategory int     `json:"max_recommendations_per_category"`
	MaxTotalRecommendations       int     `json:"max_total_recommendations"`
	MinRecommendationScore        float64 `json:"min_recommendation_score"`
	DriftThresholdForRebalance    float64 `json:"drift_threshold_for_rebalance"`
	FundingGapThreshold           float64 `json:"funding_gap_threshold"`
	DebtInterestThreshold         float64 `json:"debt_interest_threshold"`
	EmergencyFundTargetMonths     float64 `json:"emergency_fund_target_months"`
}

// DefaultConfig returns default configuration.
func DefaultConfig() Config {
	return Config{
		MaxRecommendationsPerCategory: 3,
		MaxTotalRecommendations:       20,
		MinRecommendationScore:        0.3,
		DriftThresholdForRebalance:    5.0,
		FundingGapThreshold:           100000,
		DebtInterestThreshold:         12.0,
		EmergencyFundTargetMonths:     6.0,
	}
}
