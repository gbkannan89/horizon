package domain

import (
	"errors"
	"time"
)

func ActivateLiability(l *Liability) error {
	if l.status != LSPlanned && l.status != LSApproved { return errors.New("only planned/approved liabilities can be activated") }
	l.status = LSActive; l.updatedAt = time.Now().UTC(); return nil
}

func ApplyPayment(l *Liability, amount int64) error {
	if l.status != LSActive && l.status != LSGracePeriod { return errors.New("only active liabilities accept payments") }
	if amount > l.outstandingBalance { amount = l.outstandingBalance }
	l.outstandingBalance -= amount
	if l.remainingInstallments > 0 { l.remainingInstallments-- }
	l.updatedAt = time.Now().UTC(); return nil
}

func ApplyPrepayment(l *Liability, amount int64) error {
	if l.status != LSActive { return errors.New("only active liabilities accept prepayments") }
	if amount > l.outstandingBalance { amount = l.outstandingBalance }
	l.outstandingBalance -= amount
	l.updatedAt = time.Now().UTC(); return nil
}

func EnterGracePeriod(l *Liability) error {
	if l.status != LSActive { return errors.New("only active liabilities can enter grace period") }
	l.status = LSGracePeriod; l.liabilityHealth = LHWarning; l.updatedAt = time.Now().UTC(); return nil
}

func MarkDelinquent(l *Liability) error {
	if l.status != LSActive && l.status != LSGracePeriod { return errors.New("only active/grace period liabilities can be delinquent") }
	l.status = LSDelinquent; l.liabilityHealth = LHCritical; l.updatedAt = time.Now().UTC(); return nil
}

func SettleLiability(l *Liability) error {
	if l.status != LSActive { return errors.New("only active liabilities can be settled") }
	l.status = LSSettled; l.outstandingBalance = 0; l.remainingInstallments = 0
	l.liabilityHealth = LHSettled; l.updatedAt = time.Now().UTC(); return nil
}

func WriteOffLiability(l *Liability) error {
	if l.status != LSDelinquent { return errors.New("only delinquent liabilities can be written off") }
	l.status = LSWrittenOff; l.liabilityHealth = LHHistorical; l.updatedAt = time.Now().UTC(); return nil
}

func ArchiveLiability(l *Liability) error {
	if l.status != LSSettled && l.status != LSWrittenOff && l.status != LSHistorical {
		return errors.New("only settled/written-off/historical liabilities can be archived")
	}
	l.status = LSArchived; l.updatedAt = time.Now().UTC(); return nil
}
