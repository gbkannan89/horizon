package domain

import (
	"fmt"
	"time"
)

type Goal struct {
	id                   string
	userID               string
	name                 string
	importance           GoalImportance
	gtype                GoalType
	subtype              GoalSubtype
	successCriteria      SuccessCriteria
	targetDate           *time.Time
	priority             int
	currentValue         float64
	contributionSchedule *ContributionSchedule
	riskTolerance        *RiskTolerance
	notes                string
	parentGoalID         *string
	tags                 []string
	status               GoalStatus
	createdAt            time.Time
	updatedAt            time.Time
}

func (g *Goal) ID() string                             { return g.id }
func (g *Goal) UserID() string                         { return g.userID }
func (g *Goal) Name() string                           { return g.name }
func (g *Goal) Importance() GoalImportance              { return g.importance }
func (g *Goal) GoalType() GoalType                      { return g.gtype }
func (g *Goal) Subtype() GoalSubtype                    { return g.subtype }
func (g *Goal) SuccessCriteria() SuccessCriteria        { return g.successCriteria }
func (g *Goal) TargetDate() *time.Time                  { return g.targetDate }
func (g *Goal) Priority() int                           { return g.priority }
func (g *Goal) CurrentValue() float64                   { return g.currentValue }
func (g *Goal) ContributionSchedule() *ContributionSchedule { return g.contributionSchedule }
func (g *Goal) RiskTolerance() *RiskTolerance           { return g.riskTolerance }
func (g *Goal) Notes() string                           { return g.notes }
func (g *Goal) ParentGoalID() *string                   { return g.parentGoalID }
func (g *Goal) Tags() []string                          { return g.tags }
func (g *Goal) Status() GoalStatus                      { return g.status }
func (g *Goal) CreatedAt() time.Time                    { return g.createdAt }
func (g *Goal) UpdatedAt() time.Time                    { return g.updatedAt }

func (g *Goal) SetStatus(s GoalStatus) {
	g.status = s
	g.updatedAt = time.Now().UTC()
}

func (g *Goal) SetCurrentValue(v float64) {
	g.currentValue = v
	g.updatedAt = time.Now().UTC()
}

func (g *Goal) SetPriority(p int) {
	g.priority = p
	g.updatedAt = time.Now().UTC()
}

func (g *Goal) SetTargetDate(t *time.Time) {
	g.targetDate = t
	g.updatedAt = time.Now().UTC()
}

func (g *Goal) SetSuccessCriteria(sc SuccessCriteria) {
	g.successCriteria = sc
	g.updatedAt = time.Now().UTC()
}

func (g *Goal) CanTransitionTo(target GoalStatus) error {
	transitions := map[GoalStatus]map[GoalStatus]bool{
		StatusDraft:        {StatusActive: true},
		StatusActive:       {StatusPaused: true, StatusCompleted: true, StatusMissed: true},
		StatusPaused:       {StatusActive: true, StatusArchived: true},
		StatusCompleted:    {StatusArchived: true},
		StatusMissed:       {StatusRestructured: true, StatusArchived: true},
		StatusRestructured: {StatusActive: true},
	}
	if allowed, ok := transitions[g.status]; ok {
		if allowed[target] {
			return nil
		}
	}
	if g.status == target {
		return fmt.Errorf("goal is already in %s state", target)
	}
	if g.status == StatusArchived {
		return fmt.Errorf("archived goals cannot be reactivated")
	}
	return fmt.Errorf("cannot transition from %s to %s", g.status, target)
}
