package domain

import "errors"

var (
	ErrEmptyName   = errors.New("rule name cannot be empty")
	ErrNoCondition = errors.New("rule must have at least one condition")
	ErrNoAction    = errors.New("rule must have at least one action")
)

func ValidateRule(r *Rule) error {
	if r.name == "" {
		return ErrEmptyName
	}
	if len(r.conditions) == 0 {
		return ErrNoCondition
	}
	if len(r.actions) == 0 {
		return ErrNoAction
	}
	return nil
}
