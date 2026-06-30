package mapper

import (
	"fmt"
	"time"

	"github.com/horizon/core/services/experiences/timeline/internal/engine"
)

// EventMapper converts domain/engine event types into timeline RawItems.
// This is a pure function layer — no IO, no business rules.
type EventMapper struct{}

func New() *EventMapper { return &EventMapper{} }

// MapFinancial creates a financial event timeline item.
func (m *EventMapper) MapFinancial(id, evtType, desc string, amount int64, ts time.Time) engine.RawItem {
	sev := engine.SevInfo
	if amount < 0 { sev = engine.SevWarning }
	return engine.RawItem{
		TimelineID: "fin-" + id, Timestamp: ts.Format(time.RFC3339),
		EventType: evtType, Category: engine.ECFinancial,
		Title: fmt.Sprintf("%s: %s", evtType, desc),
		Summary: fmt.Sprintf("₹%d", amount),
		Description: desc, Severity: sev,
		RelatedEntity: id, Amount: amount,
		Icon: "currency_rupee", Color: "green",
	}
}

// MapGoal creates a goal event timeline item.
func (m *EventMapper) MapGoal(id, evtType, name string, ts time.Time) engine.RawItem {
	sev, icon, color := engine.SevInfo, "track_changes", "blue"
	switch evtType {
	case "GoalCompleted": sev, icon, color = engine.SevSuccess, "celebration", "green"
	case "GoalAtRisk": sev, icon, color = engine.SevWarning, "warning", "amber"
	}
	return engine.RawItem{
		TimelineID: "goal-" + id, Timestamp: ts.Format(time.RFC3339),
		EventType: evtType, Category: engine.ECGoal,
		Title: fmt.Sprintf("%s: %s", evtType, name),
		Summary: fmt.Sprintf("Goal %s", evtType),
		RelatedEntity: id, RelatedGoal: id,
		Severity: sev, Icon: icon, Color: color,
	}
}

// MapAccount creates an account event timeline item.
func (m *EventMapper) MapAccount(id, evtType, acctType string, balance int64, ts time.Time) engine.RawItem {
	return engine.RawItem{
		TimelineID: "acct-" + id, Timestamp: ts.Format(time.RFC3339),
		EventType: evtType, Category: engine.ECAccount,
		Title: fmt.Sprintf("%s: %s", evtType, acctType),
		Summary: fmt.Sprintf("Balance: ₹%d", balance),
		RelatedEntity: id, RelatedAcct: id,
		Amount: balance, Severity: engine.SevInfo,
		Icon: "account_balance", Color: "indigo",
	}
}

// MapAsset creates an asset event timeline item.
func (m *EventMapper) MapAsset(id, evtType string, value int64, ts time.Time) engine.RawItem {
	return engine.RawItem{
		TimelineID: "ast-" + id, Timestamp: ts.Format(time.RFC3339),
		EventType: evtType, Category: engine.ECAsset,
		Title: fmt.Sprintf("Asset %s", evtType),
		Summary: fmt.Sprintf("Value: ₹%d", value),
		RelatedEntity: id, RelatedAsset: id,
		Amount: value, Severity: engine.SevInfo,
		Icon: "trending_up", Color: "teal",
	}
}

// MapLiability creates a liability event timeline item.
func (m *EventMapper) MapLiability(id, evtType string, amount int64, ts time.Time) engine.RawItem {
	sev := engine.SevWarning
	if evtType == "LiabilitySettled" { sev = engine.SevSuccess }
	if evtType == "LiabilityWrittenOff" { sev = engine.SevCritical }
	return engine.RawItem{
		TimelineID: "liab-" + id, Timestamp: ts.Format(time.RFC3339),
		EventType: evtType, Category: engine.ECLiability,
		Title: fmt.Sprintf("Liability %s", evtType),
		Summary: fmt.Sprintf("₹%d", amount),
		RelatedEntity: id, Amount: amount,
		Severity: sev, Icon: "credit_card", Color: "purple",
	}
}

// MapPortfolio creates a portfolio event timeline item.
func (m *EventMapper) MapPortfolio(id, evtType string, value int64, ts time.Time) engine.RawItem {
	return engine.RawItem{
		TimelineID: "pf-" + id, Timestamp: ts.Format(time.RFC3339),
		EventType: evtType, Category: engine.ECPortfolio,
		Title: fmt.Sprintf("Portfolio %s", evtType),
		Summary: fmt.Sprintf("Value: ₹%d", value),
		RelatedEntity: id, RelatedAgg: id,
		Amount: value, Severity: engine.SevInfo,
		Icon: "pie_chart", Color: "blue",
	}
}

// MapHealth creates a health score event timeline item.
func (m *EventMapper) MapHealth(id string, score int, grade string, ts time.Time) engine.RawItem {
	sev, icon, color := engine.SevInfo, "favorite", "pink"
	if score < 40 { sev, icon, color = engine.SevCritical, "heart_broken", "red" } else if score < 60 { sev, icon, color = engine.SevWarning, "heart_broken", "amber" }
	return engine.RawItem{
		TimelineID: "hlt-" + id, Timestamp: ts.Format(time.RFC3339),
		EventType: "HealthScoreUpdated", Category: engine.ECHealth,
		Title: fmt.Sprintf("Health Score: %d (%s)", score, grade),
		Summary: fmt.Sprintf("Score: %d/100", score),
		RelatedEntity: id, Severity: sev,
		Metadata: map[string]interface{}{"score": score, "grade": grade},
		Icon: icon, Color: color,
	}
}

// MapRisk creates a risk assessment timeline item.
func (m *EventMapper) MapRisk(id string, score int, level string, ts time.Time) engine.RawItem {
	sev, icon, color := engine.SevInfo, "shield", "green"
	if score >= 80 { sev, icon, color = engine.SevCritical, "gpp_bad", "red" } else if score >= 60 { sev, icon, color = engine.SevWarning, "shield", "amber" }
	return engine.RawItem{
		TimelineID: "risk-" + id, Timestamp: ts.Format(time.RFC3339),
		EventType: "RiskAssessmentGenerated", Category: engine.ECRisk,
		Title: fmt.Sprintf("Risk Score: %d (%s)", score, level),
		Summary: fmt.Sprintf("Risk: %s", level),
		RelatedEntity: id, Severity: sev,
		Metadata: map[string]interface{}{"score": score, "level": level},
		Icon: icon, Color: color,
	}
}

// MapRecommendation creates a recommendation timeline item.
func (m *EventMapper) MapRecommendation(id string, count int, ts time.Time) engine.RawItem {
	return engine.RawItem{
		TimelineID: "rec-" + id, Timestamp: ts.Format(time.RFC3339),
		EventType: "RecommendationsGenerated", Category: engine.ECRecommend,
		Title: fmt.Sprintf("%d New Recommendation(s)", count),
		Summary: "Review your top recommendations",
		RelatedEntity: id, Severity: engine.SevInfo,
		Metadata: map[string]interface{}{"count": count},
		Icon: "lightbulb", Color: "amber",
	}
}

// MapSimulation creates a simulation timeline item.
func (m *EventMapper) MapSimulation(id, evtType, simType string, ts time.Time) engine.RawItem {
	return engine.RawItem{
		TimelineID: "sim-" + id, Timestamp: ts.Format(time.RFC3339),
		EventType: evtType, Category: engine.ECSimulation,
		Title: fmt.Sprintf("Simulation %s: %s", simType, evtType),
		Summary: "Scenario evaluated",
		RelatedEntity: id, Severity: engine.SevInfo,
		Icon: "science", Color: "cyan",
	}
}

// MapOptimization creates an optimization timeline item.
func (m *EventMapper) MapOptimization(id, evtType string, score float64, ts time.Time) engine.RawItem {
	return engine.RawItem{
		TimelineID: "opt-" + id, Timestamp: ts.Format(time.RFC3339),
		EventType: evtType, Category: engine.ECOptimization,
		Title: fmt.Sprintf("Optimization %s", evtType),
		Summary: fmt.Sprintf("Score: %.1f", score),
		RelatedEntity: id, Severity: engine.SevInfo,
		Icon: "auto_graph", Color: "cyan",
	}
}

// MapMilestone creates a milestone timeline item.
func (m *EventMapper) MapMilestone(id, title, desc string, ts time.Time) engine.RawItem {
	return engine.RawItem{
		TimelineID: "ms-" + id, Timestamp: ts.Format(time.RFC3339),
		EventType: "MilestoneReached", Category: engine.ECMilestone,
		Title: title, Summary: desc,
		RelatedEntity: id, Severity: engine.SevMilestone,
		Icon: "flag", Color: "green",
	}
}

// MapAchievement creates an achievement timeline item.
func (m *EventMapper) MapAchievement(id, title, desc string, ts time.Time) engine.RawItem {
	return engine.RawItem{
		TimelineID: "ach-" + id, Timestamp: ts.Format(time.RFC3339),
		EventType: "AchievementUnlocked", Category: engine.ECAchievement,
		Title: title, Summary: desc,
		RelatedEntity: id, Severity: engine.SevSuccess,
		Icon: "emoji_events", Color: "gold",
	}
}

// MapUser creates a user event timeline item.
func (m *EventMapper) MapUser(id, evtType, desc string, ts time.Time) engine.RawItem {
	return engine.RawItem{
		TimelineID: "usr-" + id, Timestamp: ts.Format(time.RFC3339),
		EventType: evtType, Category: engine.ECUser,
		Title: fmt.Sprintf("%s", evtType),
		Summary: desc, RelatedEntity: id,
		Severity: engine.SevInfo, Icon: "person", Color: "grey",
	}
}
