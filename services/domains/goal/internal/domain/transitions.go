package domain

import (
	"errors"
	"time"
)

var ErrCompletedGoalImmutable = errors.New("completed goal's success criteria and target date are immutable")

func ActivateGoal(g *Goal) error {
	if g.status != StatusDraft {
		return errors.New("only draft goals can be activated")
	}
	g.status = StatusActive
	g.updatedAt = time.Now().UTC()
	return nil
}

func PauseGoal(g *Goal) error {
	if g.status != StatusActive {
		return errors.New("only active goals can be paused")
	}
	g.status = StatusPaused
	g.updatedAt = time.Now().UTC()
	return nil
}

func ResumeGoal(g *Goal) error {
	if g.status != StatusPaused {
		return errors.New("only paused goals can be resumed")
	}
	g.status = StatusActive
	g.updatedAt = time.Now().UTC()
	return nil
}

func CompleteGoal(g *Goal) error {
	if g.status != StatusActive {
		return errors.New("only active goals can be completed")
	}
	g.status = StatusCompleted
	g.updatedAt = time.Now().UTC()
	return nil
}

func MissGoal(g *Goal) error {
	if g.status != StatusActive {
		return errors.New("only active goals can be missed")
	}
	g.status = StatusMissed
	g.updatedAt = time.Now().UTC()
	return nil
}

func RestructureGoal(g *Goal, sc SuccessCriteria, targetDate *time.Time, priority int, now time.Time) error {
	if g.status != StatusActive && g.status != StatusMissed {
		return errors.New("only active or missed goals can be restructured")
	}
	if err := ValidateSuccessCriteria(sc); err != nil {
		return err
	}
	if targetDate != nil {
		if err := ValidateTargetDate(*targetDate, now); err != nil {
			return err
		}
	}
	if err := ValidatePriority(priority); err != nil {
		return err
	}
	oldSC := g.successCriteria
	oldTD := g.targetDate
	_ = oldSC
	_ = oldTD

	g.successCriteria = sc
	g.targetDate = targetDate
	g.priority = priority
	g.status = StatusRestructured
	g.updatedAt = time.Now().UTC()
	return nil
}

func CompleteRestructure(g *Goal) error {
	if g.status != StatusRestructured {
		return errors.New("only restructured goals can complete restructure")
	}
	g.status = StatusActive
	g.updatedAt = time.Now().UTC()
	return nil
}

func ArchiveGoal(g *Goal) error {
	if g.status != StatusPaused && g.status != StatusCompleted && g.status != StatusMissed {
		return errors.New("only paused, completed, or missed goals can be archived")
	}
	g.status = StatusArchived
	g.updatedAt = time.Now().UTC()
	return nil
}

func UpdatePriority(g *Goal, newPriority int) {
	g.priority = newPriority
	g.updatedAt = time.Now().UTC()
}
