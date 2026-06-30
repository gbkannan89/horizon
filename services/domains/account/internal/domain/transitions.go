package domain

import (
	"fmt"
	"time"
)

func ActivateAccount(a *Account) error {
	if a.status != StatusDraft { return fmt.Errorf("only draft accounts can be activated") }
	a.status = StatusActive; a.updatedAt = time.Now().UTC(); return nil
}
func FreezeAccount(a *Account) error {
	if a.status != StatusActive && a.status != StatusDormant {
		return fmt.Errorf("only active or dormant accounts can be frozen")
	}
	a.status = StatusFrozen; a.updatedAt = time.Now().UTC(); return nil
}
func UnfreezeAccount(a *Account) error {
	if a.status != StatusFrozen { return fmt.Errorf("only frozen accounts can be unfrozen") }
	a.status = StatusActive; a.updatedAt = time.Now().UTC(); return nil
}
func CloseAccount(a *Account, closedDate time.Time) error {
	if a.status != StatusActive && a.status != StatusFrozen {
		return fmt.Errorf("only active or frozen accounts can be closed")
	}
	a.status = StatusClosed; a.closedDate = &closedDate; a.updatedAt = time.Now().UTC(); return nil
}
func ArchiveAccount(a *Account) error {
	if a.status != StatusClosed { return fmt.Errorf("only closed accounts can be archived") }
	a.status = StatusArchived; a.updatedAt = time.Now().UTC(); return nil
}
