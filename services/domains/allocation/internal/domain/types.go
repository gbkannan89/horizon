package domain

type AllocationStatus string; const (
	StDraft AllocationStatus = "Draft"; StPlanned AllocationStatus = "Planned"
	StActive AllocationStatus = "Active"; StPaused AllocationStatus = "Paused"
	StCompleted AllocationStatus = "Completed"; StCancelled AllocationStatus = "Cancelled"
	StArchived AllocationStatus = "Archived"
)
type AllocationType string; const (
	ATPercentage AllocationType = "Percentage"; ATFixedAmount AllocationType = "FixedAmount"
	ATPriorityBased AllocationType = "PriorityBased"; ATRemainingBalance AllocationType = "RemainingBalance"
	ATRuleBased AllocationType = "RuleBased"; ATScheduled AllocationType = "Scheduled"
	ATConditional AllocationType = "Conditional"; ATEmergencyOverride AllocationType = "EmergencyOverride"
	ATManual AllocationType = "Manual"; ATHybrid AllocationType = "Hybrid"
)
type AllocationStrategy string; const (
	ASEqual AllocationStrategy = "Equal"; ASWeighted AllocationStrategy = "Weighted"
	ASGoalPriority AllocationStrategy = "GoalPriority"; ASGoalImportance AllocationStrategy = "GoalImportance"
	ASTargetDate AllocationStrategy = "TargetDate"; ASRiskBased AllocationStrategy = "RiskBased"
	ASCashFlowBased AllocationStrategy = "CashFlowBased"; ASCustom AllocationStrategy = "Custom"
	ASAISuggested AllocationStrategy = "AISuggested"
)
type AllocationHealth string; const (
	AHHealthy AllocationHealth = "Healthy"; AHUnderfunded AllocationHealth = "Underfunded"
	AHAtRisk AllocationHealth = "AtRisk"; AHOverallocated AllocationHealth = "Overallocated"
	AHBlocked AllocationHealth = "Blocked"; AHCompletedAlloc AllocationHealth = "Completed"
)
type AllocationConfidence string; const (
	ACConfirmed AllocationConfidence = "Confirmed"; ACPlanned AllocationConfidence = "Planned"
	ACEstimated AllocationConfidence = "Estimated"; ACSimulated AllocationConfidence = "Simulated"
	ACInferred AllocationConfidence = "Inferred"
)
type FundingSourceType string; const (
	FSTAccount FundingSourceType = "Account"; FSTFinancialEvent FundingSourceType = "FinancialEvent"
	FSTVirtualEnvelope FundingSourceType = "VirtualEnvelope"; FSTIncomeStream FundingSourceType = "IncomeStream"
)
