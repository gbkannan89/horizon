package domain

import (
	"fmt"
	"time"
)

func ActivateBudget(b *Budget) error {
	if b.status != BStatusDraft {
		return fmt.Errorf("only draft budgets can be activated")
	}
	b.status = BStatusActive
	b.updatedAt = time.Now().UTC()
	return nil
}

func PauseBudget(b *Budget) error {
	if b.status != BStatusActive {
		return fmt.Errorf("only active budgets can be paused")
	}
	b.status = BStatusPaused
	b.updatedAt = time.Now().UTC()
	return nil
}

func ResumeBudget(b *Budget) error {
	if b.status != BStatusPaused {
		return fmt.Errorf("only paused budgets can be resumed")
	}
	b.status = BStatusActive
	b.updatedAt = time.Now().UTC()
	return nil
}

func CompleteBudget(b *Budget) error {
	if b.status != BStatusActive {
		return fmt.Errorf("only active budgets can be completed")
	}
	b.status = BStatusCompleted
	b.updatedAt = time.Now().UTC()
	return nil
}

func ArchiveBudget(b *Budget) error {
	if b.status != BStatusCompleted && b.status != BStatusPaused {
		return fmt.Errorf("only completed or paused budgets can be archived")
	}
	b.status = BStatusArchived
	b.updatedAt = time.Now().UTC()
	return nil
}
