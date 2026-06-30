package domain

import (
	"errors"
	"time"
)

type GoalFactory struct{}

func NewGoalFactory() *GoalFactory { return &GoalFactory{} }

func (f *GoalFactory) Create(
	id, userID, name string,
	importance GoalImportance,
	gtype GoalType,
	subtype GoalSubtype,
	sc SuccessCriteria,
	priority int,
	targetDate *time.Time,
	targetDateValid bool,
	riskTolerance string,
	notes string,
	parentGoalID *string,
	tags []string,
	schedule *ContributionSchedule,
	now time.Time,
) (*Goal, error) {

	if err := ValidateGoalName(name); err != nil {
		return nil, err
	}
	if _, ok := AllImportance[importance]; !ok {
		return nil, errors.New("importance must be Mandatory, Essential, Lifestyle, Dream, or Speculative")
	}
	if err := ValidateSuccessCriteria(sc); err != nil {
		return nil, err
	}
	if err := ValidatePriority(priority); err != nil {
		return nil, err
	}
	if riskTolerance != "" {
		if err := ValidateRiskTolerance(riskTolerance); err != nil {
			return nil, err
		}
	}
	if gtype == TypeTimeBound && targetDate == nil {
		return nil, errors.New("target date is required for time-bound goals")
	}
	if targetDate != nil {
		if err := ValidateTargetDate(*targetDate, now); err != nil {
			return nil, err
		}
	}

	return &Goal{
		id:                   id,
		userID:               userID,
		name:                 name,
		importance:           importance,
		gtype:                gtype,
		subtype:              subtype,
		successCriteria:      sc,
		priority:             priority,
		targetDate:           targetDate,
		currentValue:         0,
		contributionSchedule: schedule,
		riskTolerance:        getRiskPtr(riskTolerance),
		notes:                notes,
		parentGoalID:         parentGoalID,
		tags:                 tags,
		status:               StatusDraft,
		createdAt:            now,
		updatedAt:            now,
	}, nil
}

func getRiskPtr(s string) *RiskTolerance {
	if s == "" {
		return nil
	}
	r := RiskTolerance(s)
	return &r
}
