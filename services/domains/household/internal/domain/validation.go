package domain

import (
	"errors"
	"strings"
)

func ValidateHouseholdName(n string) error {
	t := strings.TrimSpace(n)
	if t == "" { return errors.New("household name is required") }
	if len([]rune(t)) > 200 { return errors.New("household name must be 200 characters or fewer") }
	return nil
}

func ValidateHouseholdType(t HouseholdType) error {
	switch t {
	case HHSingle, HHCouple, HHFamily, HHRoommates, HHCustom:
		return nil
	}
	return errors.New("household type is not recognized")
}

func ValidateMemberRole(r MemberRole) error {
	switch r {
	case RoleHead, RoleAdmin, RoleMember, RoleViewer:
		return nil
	}
	return errors.New("member role is not recognized")
}

func ValidateCurrency(c string) error {
	if len(c) != 3 { return errors.New("currency must be a 3-letter ISO code") }
	return nil
}

func ValidateMemberCount(members []HouseholdMember) error {
	if len(members) == 0 { return errors.New("household must have at least one member") }
	if len(members) > 50 { return errors.New("household cannot have more than 50 members") }
	return nil
}

func ValidateHeadOfHousehold(members []HouseholdMember, headID string) error {
	found := false
	for _, m := range members {
		if m.UserID == headID {
			found = true
			if m.Role != RoleHead {
				return errors.New("head of household must have Head role")
			}
		}
	}
	if !found { return errors.New("head of household must be a member") }
	return nil
}
