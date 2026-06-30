package domain

import "time"

type DomainEvent interface {
	EventName() string
	EntityID() string
}

type Publisher interface {
	Publish(event DomainEvent) error
}

type BaseDomainEvent struct {
	Type string    `json:"event_type"`
	ID   string    `json:"entity_id"`
	Time time.Time `json:"timestamp"`
}

func NewBaseDomainEvent(eventType, entityID string) BaseDomainEvent {
	return BaseDomainEvent{
		Type: eventType,
		ID:   entityID,
		Time: time.Now().UTC(),
	}
}

func (e BaseDomainEvent) EventName() string { return e.Type }
func (e BaseDomainEvent) EntityID() string  { return e.ID }

type UserRegistered struct {
	BaseDomainEvent
	UserID      string `json:"user_id"`
	UserType    string `json:"user_type"`
	Country     string `json:"country"`
	BaseCurrency string `json:"base_currency"`
}

func NewUserRegistered(u *User) UserRegistered {
	return UserRegistered{
		BaseDomainEvent: NewBaseDomainEvent("UserRegistered", u.UserID()),
		UserID:          u.UserID(),
		UserType:        string(u.UserType()),
		Country:         u.Country(),
		BaseCurrency:    u.BaseCurrency(),
	}
}

type UserActivated struct {
	BaseDomainEvent
	UserID string `json:"user_id"`
}

func NewUserActivated(u *User) UserActivated {
	return UserActivated{
		BaseDomainEvent: NewBaseDomainEvent("UserActivated", u.UserID()),
		UserID:          u.UserID(),
	}
}

type UserSuspended struct {
	BaseDomainEvent
	UserID string `json:"user_id"`
	Reason string `json:"reason"`
}

func NewUserSuspended(u *User, reason string) UserSuspended {
	return UserSuspended{
		BaseDomainEvent: NewBaseDomainEvent("UserSuspended", u.UserID()),
		UserID:          u.UserID(),
		Reason:          reason,
	}
}

type UserArchived struct {
	BaseDomainEvent
	UserID string `json:"user_id"`
	Reason string `json:"reason,omitempty"`
}

func NewUserArchived(u *User, reason string) UserArchived {
	return UserArchived{
		BaseDomainEvent: NewBaseDomainEvent("UserArchived", u.UserID()),
		UserID:          u.UserID(),
		Reason:          reason,
	}
}

type ProfileUpdated struct {
	BaseDomainEvent
	UserID       string                 `json:"user_id"`
	ChangedFields []string              `json:"changed_fields"`
	OldValues    map[string]interface{} `json:"old_values"`
	NewValues    map[string]interface{} `json:"new_values"`
}

func NewProfileUpdated(u *User, changed []string, old, new map[string]interface{}) ProfileUpdated {
	return ProfileUpdated{
		BaseDomainEvent: NewBaseDomainEvent("ProfileUpdated", u.UserID()),
		UserID:          u.UserID(),
		ChangedFields:   changed,
		OldValues:       old,
		NewValues:       new,
	}
}

type PreferencesUpdated struct {
	BaseDomainEvent
	UserID            string                 `json:"user_id"`
	ChangedPreferences map[string]interface{} `json:"changed_preferences"`
}

func NewPreferencesUpdated(u *User, changed map[string]interface{}) PreferencesUpdated {
	return PreferencesUpdated{
		BaseDomainEvent:    NewBaseDomainEvent("PreferencesUpdated", u.UserID()),
		UserID:             u.UserID(),
		ChangedPreferences: changed,
	}
}

type ConsentUpdated struct {
	BaseDomainEvent
	UserID      string `json:"user_id"`
	ConsentType string `json:"consent_type"`
	Action      string `json:"action"`
}

func NewConsentUpdated(u *User, cType ConsentType, action string) ConsentUpdated {
	return ConsentUpdated{
		BaseDomainEvent: NewBaseDomainEvent("ConsentUpdated", u.UserID()),
		UserID:          u.UserID(),
		ConsentType:     string(cType),
		Action:          action,
	}
}

type PrivacyUpdated struct {
	BaseDomainEvent
	UserID         string         `json:"user_id"`
	ChangedSettings map[string]interface{} `json:"changed_settings"`
}

func NewPrivacyUpdated(u *User, changed map[string]interface{}) PrivacyUpdated {
	return PrivacyUpdated{
		BaseDomainEvent:  NewBaseDomainEvent("PrivacyUpdated", u.UserID()),
		UserID:           u.UserID(),
		ChangedSettings:  changed,
	}
}

type HouseholdJoined struct {
	BaseDomainEvent
	UserID      string `json:"user_id"`
	HouseholdID string `json:"household_id"`
	Role        string `json:"role"`
}

type HouseholdLeft struct {
	BaseDomainEvent
	UserID      string `json:"user_id"`
	HouseholdID string `json:"household_id"`
}
