package domain

import "time"

type DomainEvent interface {
	EventName() string
	EntityID() string
}

type Publisher interface {
	Publish(event DomainEvent) error
}

type BaseDomainEvent struct {
	Type string    `json:"event_type"`
	ID   string    `json:"entity_id"`
	Time time.Time `json:"timestamp"`
}

func NewBaseDomainEvent(eventType, entityID string) BaseDomainEvent {
	return BaseDomainEvent{Type: eventType, ID: entityID, Time: time.Now().UTC()}
}

func (e BaseDomainEvent) EventName() string { return e.Type }
func (e BaseDomainEvent) EntityID() string  { return e.ID }

type GoalCreated struct {
	BaseDomainEvent
	UserID         string          `json:"user_id"`
	Name           string          `json:"name"`
	Importance     string          `json:"importance"`
	GoalType       string          `json:"goal_type"`
	Priority       int             `json:"priority"`
	SuccessCriteria SuccessCriteria `json:"success_criteria"`
}
func NewGoalCreated(g *Goal) GoalCreated {
	return GoalCreated{
		BaseDomainEvent: NewBaseDomainEvent("GoalCreated", g.ID()),
		UserID: g.UserID(), Name: g.Name(), Importance: string(g.Importance()),
		GoalType: string(g.GoalType()), Priority: g.Priority(), SuccessCriteria: g.SuccessCriteria(),
	}
}

type GoalActivated struct {
	BaseDomainEvent
	UserID         string `json:"user_id"`
	ActivationDate string `json:"activation_date"`
}
func NewGoalActivated(g *Goal) GoalActivated {
	return GoalActivated{
		BaseDomainEvent: NewBaseDomainEvent("GoalActivated", g.ID()),
		UserID: g.UserID(), ActivationDate: time.Now().UTC().Format(time.RFC3339),
	}
}

type GoalProgressed struct {
	BaseDomainEvent
	NewProgress float64 `json:"new_progress"`
	Delta       float64 `json:"delta"`
}
func NewGoalProgressed(g *Goal, delta float64) GoalProgressed {
	return GoalProgressed{
		BaseDomainEvent: NewBaseDomainEvent("GoalProgressed", g.ID()),
		NewProgress: g.CurrentValue(), Delta: delta,
	}
}

type GoalPaused struct {
	BaseDomainEvent
	Reason string `json:"reason,omitempty"`
}
func NewGoalPaused(g *Goal, reason string) GoalPaused {
	return GoalPaused{
		BaseDomainEvent: NewBaseDomainEvent("GoalPaused", g.ID()),
		Reason: reason,
	}
}

type GoalResumed struct {
	BaseDomainEvent
}
func NewGoalResumed(g *Goal) GoalResumed {
	return GoalResumed{BaseDomainEvent: NewBaseDomainEvent("GoalResumed", g.ID())}
}

type GoalRestructuredEvent struct {
	BaseDomainEvent
}
func NewGoalRestructured(g *Goal) GoalRestructuredEvent {
	return GoalRestructuredEvent{BaseDomainEvent: NewBaseDomainEvent("GoalRestructured", g.ID())}
}

type GoalCompletedEvent struct {
	BaseDomainEvent
	FinalValue     float64 `json:"final_value"`
	CompletionDate string  `json:"completion_date"`
}
func NewGoalCompleted(g *Goal) GoalCompletedEvent {
	return GoalCompletedEvent{
		BaseDomainEvent: NewBaseDomainEvent("GoalCompleted", g.ID()),
		FinalValue: g.CurrentValue(), CompletionDate: time.Now().UTC().Format(time.RFC3339),
	}
}

type GoalMissedEvent struct {
	BaseDomainEvent
	FinalValue float64 `json:"final_value"`
	MissedDate string  `json:"missed_date"`
}
func NewGoalMissed(g *Goal) GoalMissedEvent {
	return GoalMissedEvent{
		BaseDomainEvent: NewBaseDomainEvent("GoalMissed", g.ID()),
		FinalValue: g.CurrentValue(), MissedDate: time.Now().UTC().Format(time.RFC3339),
	}
}

type GoalArchivedEvent struct {
	BaseDomainEvent
}
func NewGoalArchived(g *Goal) GoalArchivedEvent {
	return GoalArchivedEvent{BaseDomainEvent: NewBaseDomainEvent("GoalArchived", g.ID())}
}

type GoalPriorityChanged struct {
	BaseDomainEvent
	OldPriority int `json:"old_priority"`
	NewPriority int `json:"new_priority"`
}
func NewGoalPriorityChanged(g *Goal, old int) GoalPriorityChanged {
	return GoalPriorityChanged{
		BaseDomainEvent: NewBaseDomainEvent("GoalPriorityChanged", g.ID()),
		OldPriority: old, NewPriority: g.Priority(),
	}
}

type GoalAtRisk struct {
	BaseDomainEvent
	DeviationMagnitude float64 `json:"deviation_magnitude"`
	ProjectedShortfall  float64 `json:"projected_shortfall"`
}
func NewGoalAtRisk(g *Goal, deviation, shortfall float64) GoalAtRisk {
	return GoalAtRisk{
		BaseDomainEvent: NewBaseDomainEvent("GoalAtRisk", g.ID()),
		DeviationMagnitude: deviation, ProjectedShortfall: shortfall,
	}
}

type GoalBackOnTrack struct {
	BaseDomainEvent
}
func NewGoalBackOnTrack(g *Goal) GoalBackOnTrack {
	return GoalBackOnTrack{BaseDomainEvent: NewBaseDomainEvent("GoalBackOnTrack", g.ID())}
}
