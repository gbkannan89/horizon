package domain

import (
	"errors"
	"strings"
	"time"
)

func ValidateName(n string) error { t := strings.TrimSpace(n); if t == "" { return errors.New("name is required") }; if len([]rune(t)) > 300 { return errors.New("name must be 300 characters or fewer") }; return nil }
func ValidateCountry(c string) error { if c == "" { return errors.New("country is required") }; return nil }

func RegisterInstitution(i *Institution) error {
	if i.status != StDraft { return errors.New("only draft institutions can be registered") }
	i.status = StRegistered; i.trustLevel = TLUnverified; i.health = IHStable; i.confidence = ICRegistered; i.updatedAt = time.Now().UTC(); return nil
}
func VerifyInstitution(i *Institution) error {
	if i.status != StRegistered { return errors.New("only registered institutions can be verified") }
	i.status = StVerified; i.trustLevel = TLVerified; i.confidence = ICVerified; i.updatedAt = time.Now().UTC(); return nil
}
func ActivateInstitution(i *Institution) error {
	if i.status != StVerified { return errors.New("only verified institutions can be activated") }
	i.status = StActive; i.health = IHHealthy; i.updatedAt = time.Now().UTC(); return nil
}
func SuspendInstitution(i *Institution) error {
	if i.status != StActive && i.status != StVerified { return errors.New("only active or verified institutions can be suspended") }
	i.status = StSuspended; i.health = IHSuspended; i.updatedAt = time.Now().UTC(); return nil
}
func CloseInstitution(i *Institution) error {
	if i.status != StActive && i.status != StSuspended { return errors.New("only active or suspended institutions can be closed") }
	i.status = StClosed; i.health = IHHistorical; i.updatedAt = time.Now().UTC(); return nil
}
func ArchiveInstitution(i *Institution) error {
	if i.status != StClosed && i.status != StMerged { return errors.New("only closed or merged institutions can be archived") }
	i.status = StArchived; i.updatedAt = time.Now().UTC(); return nil
}
