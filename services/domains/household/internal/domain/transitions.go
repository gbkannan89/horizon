package domain

import (
	"fmt"
	"time"
)

func ActivateHousehold(h *Household) error {
	if h.status != HHStatusDraft {
		return fmt.Errorf("only draft households can be activated")
	}
	h.status = HHStatusActive
	h.updatedAt = time.Now().UTC()
	return nil
}

func PauseHousehold(h *Household) error {
	if h.status != HHStatusActive {
		return fmt.Errorf("only active households can be paused")
	}
	h.status = HHStatusPaused
	h.updatedAt = time.Now().UTC()
	return nil
}

func ResumeHousehold(h *Household) error {
	if h.status != HHStatusPaused {
		return fmt.Errorf("only paused households can be resumed")
	}
	h.status = HHStatusActive
	h.updatedAt = time.Now().UTC()
	return nil
}

func DissolveHousehold(h *Household) error {
	if h.status != HHStatusActive && h.status != HHStatusPaused {
		return fmt.Errorf("only active or paused households can be dissolved")
	}
	h.status = HHStatusDissolved
	h.updatedAt = time.Now().UTC()
	return nil
}

func ArchiveHousehold(h *Household) error {
	if h.status != HHStatusDissolved {
		return fmt.Errorf("only dissolved households can be archived")
	}
	h.status = HHStatusArchived
	h.updatedAt = time.Now().UTC()
	return nil
}
