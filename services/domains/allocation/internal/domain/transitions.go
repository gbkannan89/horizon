package domain

import (
	"errors"
	"time"
)

func ValidatePriority(p int) error { if p < 1 { return errors.New("priority must be >= 1") }; return nil }
func ValidateAmount(v int64) error { if v <= 0 { return errors.New("amount must be > 0") }; return nil }
func ValidatePercentage(v float64) error { if v < 0 || v > 100 { return errors.New("percentage must be 0-100") }; return nil }
func ValidateEffectiveDate(d time.Time) error { if d.Before(time.Now().Truncate(24 * time.Hour)) { return errors.New("effective date must be today or later") }; return nil }
func ValidateExpirationDate(eff, exp time.Time) error { if exp.Before(eff) || exp.Equal(eff) { return errors.New("expiration must be after effective date") }; return nil }
func ValidateCurrencyMatch(allocCurr, sourceCurr string) error { if allocCurr != sourceCurr { return errors.New("currency must match funding source") }; return nil }
func ValidateApprovalNeeded(cb string, ab *string) error { if cb == "Recommendation" && (ab == nil || *ab == "") { return errors.New("approval required for recommendation-created allocations") }; return nil }

func ActivateAlloc(a *Allocation) error {
	if a.status != StPlanned { return errors.New("only planned allocations can be activated") }
	a.status = StActive; a.updatedAt = time.Now().UTC(); return nil
}
func PauseAlloc(a *Allocation) error {
	if a.status != StActive { return errors.New("only active allocations can be paused") }
	a.status = StPaused; a.reservedAmount = 0; a.updatedAt = time.Now().UTC(); return nil
}
func ResumeAlloc(a *Allocation) error {
	if a.status != StPaused { return errors.New("only paused allocations can be resumed") }
	a.status = StActive; a.updatedAt = time.Now().UTC(); return nil
}
func CompleteAlloc(a *Allocation) error {
	if a.status != StActive { return errors.New("only active allocations can be completed") }
	a.status = StCompleted; a.reservedAmount = 0; a.updatedAt = time.Now().UTC(); return nil
}
func CancelAlloc(a *Allocation) error {
	if a.status != StDraft && a.status != StPlanned && a.status != StPaused { return errors.New("only draft/planned/paused allocations can be cancelled") }
	a.status = StCancelled; a.reservedAmount = 0; a.updatedAt = time.Now().UTC(); return nil
}
func ArchiveAlloc(a *Allocation) error {
	if a.status != StCompleted && a.status != StCancelled { return errors.New("only completed or cancelled allocations can be archived") }
	a.status = StArchived; a.updatedAt = time.Now().UTC(); return nil
}
