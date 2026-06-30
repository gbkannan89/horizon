package domain

import "errors"

type Transition func(*FinancialEvent) error

func SubmitEvent(event *FinancialEvent) error {
	if err := event.CanTransitionTo(StatePending); err != nil {
		return err
	}
	event.state = StatePending
	return nil
}

func ConfirmEvent(event *FinancialEvent) error {
	if err := event.CanTransitionTo(StateConfirmed); err != nil {
		return err
	}
	if event.state != StatePending {
		return errors.New("only pending events can be confirmed")
	}
	event.state = StateConfirmed
	return nil
}

func PostEvent(event *FinancialEvent) error {
	if err := event.CanTransitionTo(StatePosted); err != nil {
		return err
	}
	if event.state != StateConfirmed {
		return errors.New("only confirmed events can be posted")
	}
	event.state = StatePosted
	return nil
}

func CancelEvent(event *FinancialEvent) error {
	if err := event.CanTransitionTo(StateCancelled); err != nil {
		return err
	}
	if event.state != StateDraft && event.state != StatePending {
		return errors.New("only draft or pending events can be cancelled")
	}
	event.state = StateCancelled
	return nil
}

func ArchiveEvent(event *FinancialEvent) error {
	if err := event.CanTransitionTo(StateArchived); err != nil {
		return err
	}
	if event.state != StatePosted && event.state != StateReversed && event.state != StateCancelled {
		return errors.New("only posted, reversed, or cancelled events can be archived")
	}
	event.state = StateArchived
	return nil
}

func ReverseEvent(original, reversal *FinancialEvent) error {
	if err := original.CanTransitionTo(StateReversed); err != nil {
		return err
	}
	if original.state != StateConfirmed && original.state != StatePosted {
		return errors.New("only confirmed or posted events can be reversed")
	}
	if reversal.reversalOfEventID == "" {
		return errors.New("reversal event must reference the original event")
	}
	if reversal.amount != -original.amount {
		return errors.New("reversal amount must be the negation of the original amount")
	}
	original.state = StateReversed
	return nil
}

func UpgradeConfidence(event *FinancialEvent, newConfidence EventConfidence) error {
	if event.state != StateDraft && event.state != StatePending {
		return errors.New("confidence can only be changed in Draft or Pending state")
	}
	if !event.confidence.CanUpgradeTo(newConfidence) {
		return errors.New("confidence may not be downgraded")
	}
	event.confidence = newConfidence
	return nil
}
