package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/horizon/core/services/domains/user/internal/application/dto/command"
	"github.com/horizon/core/services/domains/user/internal/application/dto/query"
	"github.com/horizon/core/services/domains/user/internal/domain"
)

type UserService struct {
	repo      domain.Repository
	publisher domain.Publisher
	factory   *domain.UserFactory
	now       func() time.Time
}

func NewUserService(repo domain.Repository, publisher domain.Publisher, now func() time.Time) *UserService {
	return &UserService{
		repo:      repo,
		publisher: publisher,
		factory:   domain.NewUserFactory(),
		now:       now,
	}
}

func (s *UserService) Register(ctx context.Context, cmd command.RegisterUserCommand) (*command.RegisterResult, error) {
	userType := domain.UserType(cmd.UserType)
	if cmd.UserType == "" {
		userType = domain.TypeIndividual
	}
	fip := domain.FinancialIdentityProfile(cmd.FinancialIdentityProfile)
	if cmd.FinancialIdentityProfile == "" {
		fip = domain.FIPPersonal
	}
	sot := domain.SOT(cmd.SourceOfTruth)
	if cmd.SourceOfTruth == "" {
		sot = domain.SOTSystem
	}

	var dob *time.Time
	if cmd.DateOfBirth != "" {
		t, err := time.Parse("2006-01-02", cmd.DateOfBirth)
		if err != nil {
			return nil, fmt.Errorf("invalid date of birth: %w", err)
		}
		dob = &t
	}

	u, err := s.factory.Register(
		cmd.UserID, cmd.DisplayName, cmd.LegalName, cmd.Country, cmd.BaseCurrency,
		cmd.Locale, cmd.Timezone, userType, cmd.PreferredName, dob, fip, sot, cmd.HouseholdID,
	)
	if err != nil {
		return nil, err
	}

	u.SetProfileCompleteness(domain.CalculateProfileCompleteness(u))

	if err := s.repo.Save(ctx, u); err != nil {
		return nil, fmt.Errorf("save user: %w", err)
	}
	if err := s.publisher.Publish(domain.NewUserRegistered(u)); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	return &command.RegisterResult{UserID: u.UserID(), Status: string(u.Status())}, nil
}

func (s *UserService) Activate(ctx context.Context, cmd command.ActivateUserCommand) (*command.CommandResult, error) {
	u, err := s.repo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %s", cmd.UserID)
	}
	if err := u.CanTransitionTo(domain.StatusActive); err != nil {
		return nil, err
	}
	u.SetStatus(domain.StatusActive)
	if err := s.repo.Save(ctx, u); err != nil {
		return nil, fmt.Errorf("save user: %w", err)
	}
	if err := s.publisher.Publish(domain.NewUserActivated(u)); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	return &command.CommandResult{Success: true, UserID: u.UserID(), Status: string(u.Status())}, nil
}

func (s *UserService) Suspend(ctx context.Context, cmd command.SuspendUserCommand) (*command.CommandResult, error) {
	u, err := s.repo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %s", cmd.UserID)
	}
	if err := u.CanTransitionTo(domain.StatusSuspended); err != nil {
		return nil, err
	}
	u.SetStatus(domain.StatusSuspended)
	if err := s.repo.Save(ctx, u); err != nil {
		return nil, fmt.Errorf("save user: %w", err)
	}
	if err := s.publisher.Publish(domain.NewUserSuspended(u, cmd.Reason)); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	return &command.CommandResult{Success: true, UserID: u.UserID(), Status: string(u.Status())}, nil
}

func (s *UserService) Reactivate(ctx context.Context, cmd command.ReactivateUserCommand) (*command.CommandResult, error) {
	u, err := s.repo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %s", cmd.UserID)
	}
	if err := u.CanTransitionTo(domain.StatusActive); err != nil {
		return nil, err
	}
	if u.Status() != domain.StatusSuspended && u.Status() != domain.StatusInactive {
		return nil, errors.New("only suspended or inactive users can be reactivated")
	}
	u.SetStatus(domain.StatusActive)
	if err := s.repo.Save(ctx, u); err != nil {
		return nil, fmt.Errorf("save user: %w", err)
	}
	if err := s.publisher.Publish(domain.NewUserActivated(u)); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	return &command.CommandResult{Success: true, UserID: u.UserID(), Status: string(u.Status())}, nil
}

func (s *UserService) Archive(ctx context.Context, cmd command.ArchiveUserCommand) (*command.CommandResult, error) {
	u, err := s.repo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %s", cmd.UserID)
	}
	if err := u.CanTransitionTo(domain.StatusArchived); err != nil {
		return nil, err
	}
	u.SetStatus(domain.StatusArchived)
	if err := s.repo.Save(ctx, u); err != nil {
		return nil, fmt.Errorf("save user: %w", err)
	}
	if err := s.publisher.Publish(domain.NewUserArchived(u, cmd.Reason)); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	return &command.CommandResult{Success: true, UserID: u.UserID(), Status: string(u.Status())}, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, cmd command.UpdateProfileCommand) (*command.CommandResult, error) {
	u, err := s.repo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %s", cmd.UserID)
	}
	changed := []string{}
	oldVals := make(map[string]interface{})
	newVals := make(map[string]interface{})

	if cmd.DisplayName != "" && cmd.DisplayName != u.DisplayName() {
		if err := domain.ValidateDisplayName(cmd.DisplayName); err != nil {
			return nil, err
		}
		oldVals["display_name"] = u.DisplayName()
		u.SetDisplayName(cmd.DisplayName)
		newVals["display_name"] = cmd.DisplayName
		changed = append(changed, "display_name")
	}
	if cmd.LegalName != "" && cmd.LegalName != u.LegalName() {
		if err := domain.ValidateLegalName(cmd.LegalName); err != nil {
			return nil, err
		}
		oldVals["legal_name"] = u.LegalName()
		u.SetLegalName(cmd.LegalName)
		newVals["legal_name"] = cmd.LegalName
		changed = append(changed, "legal_name")
	}
	if cmd.Country != "" && cmd.Country != u.Country() {
		if err := domain.ValidateCountry(cmd.Country); err != nil {
			return nil, err
		}
		oldVals["country"] = u.Country()
		u.SetCountry(cmd.Country)
		newVals["country"] = cmd.Country
		changed = append(changed, "country")
	}
	if cmd.BaseCurrency != "" && cmd.BaseCurrency != u.BaseCurrency() {
		if err := domain.ValidateBaseCurrency(cmd.BaseCurrency); err != nil {
			return nil, err
		}
		oldVals["base_currency"] = u.BaseCurrency()
		u.SetBaseCurrency(cmd.BaseCurrency)
		newVals["base_currency"] = cmd.BaseCurrency
		changed = append(changed, "base_currency")
	}
	if cmd.Locale != "" && cmd.Locale != u.Locale() {
		if err := domain.ValidateLocale(cmd.Locale); err != nil {
			return nil, err
		}
		oldVals["locale"] = u.Locale()
		u.SetLocale(cmd.Locale)
		newVals["locale"] = cmd.Locale
		changed = append(changed, "locale")
	}
	if cmd.Timezone != "" && cmd.Timezone != u.Timezone() {
		if err := domain.ValidateTimezone(cmd.Timezone); err != nil {
			return nil, err
		}
		oldVals["timezone"] = u.Timezone()
		u.SetTimezone(cmd.Timezone)
		newVals["timezone"] = cmd.Timezone
		changed = append(changed, "timezone")
	}

	if len(changed) == 0 {
		return &command.CommandResult{Success: true, UserID: u.UserID(), Status: string(u.Status())}, nil
	}

	if err := s.repo.Save(ctx, u); err != nil {
		return nil, fmt.Errorf("save user: %w", err)
	}
	if err := s.publisher.Publish(domain.NewProfileUpdated(u, changed, oldVals, newVals)); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	return &command.CommandResult{Success: true, UserID: u.UserID(), Status: string(u.Status())}, nil
}

func (s *UserService) UpdatePreferences(ctx context.Context, cmd command.UpdatePreferencesCommand) (*command.CommandResult, error) {
	u, err := s.repo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %s", cmd.UserID)
	}
	u.SetPreferences(cmd.Preferences)
	if err := s.repo.Save(ctx, u); err != nil {
		return nil, fmt.Errorf("save user: %w", err)
	}
	if err := s.publisher.Publish(domain.NewPreferencesUpdated(u, cmd.Preferences)); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	return &command.CommandResult{Success: true, UserID: u.UserID(), Status: string(u.Status())}, nil
}

func (s *UserService) GrantConsent(ctx context.Context, cmd command.GrantConsentCommand) (*command.CommandResult, error) {
	u, err := s.repo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %s", cmd.UserID)
	}
	cType := domain.ConsentType(cmd.ConsentType)
	if !domain.AllConsentTypes[cType] {
		return nil, fmt.Errorf("invalid consent type: %s", cmd.ConsentType)
	}
	u.AddConsent(domain.ConsentRecord{
		ConsentType: cType,
		Granted:     true,
		GrantedAt:   s.now().UTC(),
		Scope:       cmd.Scope,
	})
	if err := s.repo.Save(ctx, u); err != nil {
		return nil, fmt.Errorf("save user: %w", err)
	}
	if err := s.publisher.Publish(domain.NewConsentUpdated(u, cType, "Granted")); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	return &command.CommandResult{Success: true, UserID: u.UserID()}, nil
}

func (s *UserService) RevokeConsent(ctx context.Context, cmd command.RevokeConsentCommand) (*command.CommandResult, error) {
	u, err := s.repo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %s", cmd.UserID)
	}
	cType := domain.ConsentType(cmd.ConsentType)
	if !u.RevokeConsent(cType) {
		return nil, fmt.Errorf("no active consent found for type: %s", cmd.ConsentType)
	}
	if err := s.repo.Save(ctx, u); err != nil {
		return nil, fmt.Errorf("save user: %w", err)
	}
	if err := s.publisher.Publish(domain.NewConsentUpdated(u, cType, "Revoked")); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	return &command.CommandResult{Success: true, UserID: u.UserID()}, nil
}

func (s *UserService) UpdatePrivacy(ctx context.Context, cmd command.UpdatePrivacyCommand) (*command.CommandResult, error) {
	u, err := s.repo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %s", cmd.UserID)
	}
	changed := make(map[string]interface{})
	pp := u.PrivacyProfile()
	if cmd.DataSharing != "" {
		pp.DataSharing = domain.DataSharingLevel(cmd.DataSharing)
		changed["data_sharing"] = cmd.DataSharing
	}
	if cmd.HouseholdVis != "" {
		pp.HouseholdVisibility = domain.DataSharingLevel(cmd.HouseholdVis)
		changed["household_visibility"] = cmd.HouseholdVis
	}
	if cmd.ThirdPartyAccess != nil {
		pp.ThirdPartyAccess = *cmd.ThirdPartyAccess
		changed["third_party_access"] = *cmd.ThirdPartyAccess
	}
	u.SetPrivacy(pp)
	if err := s.repo.Save(ctx, u); err != nil {
		return nil, fmt.Errorf("save user: %w", err)
	}
	if err := s.publisher.Publish(domain.NewPrivacyUpdated(u, changed)); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	return &command.CommandResult{Success: true, UserID: u.UserID()}, nil
}

func (s *UserService) GetUser(ctx context.Context, q query.GetUserQuery) (*query.UserResult, error) {
	u, err := s.repo.GetByID(ctx, q.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %s", q.UserID)
	}
	return toUserResult(u), nil
}

func (s *UserService) ListByStatus(ctx context.Context, q query.ListUsersByStatusQuery) (*query.PaginatedResult, error) {
	status := domain.UserStatus(q.Status)
	users, cursor, err := s.repo.ListByStatus(ctx, status, q.Cursor, q.Limit)
	if err != nil {
		return nil, err
	}
	return toPaginatedResult(users, cursor), nil
}

func (s *UserService) Search(ctx context.Context, q query.SearchUsersQuery) (*query.PaginatedResult, error) {
	users, cursor, err := s.repo.Search(ctx, q.Query, q.Cursor, q.Limit)
	if err != nil {
		return nil, err
	}
	return toPaginatedResult(users, cursor), nil
}

func toUserResult(u *domain.User) *query.UserResult {
	return &query.UserResult{
		UserID:                   u.UserID(),
		DisplayName:              u.DisplayName(),
		LegalName:                u.LegalName(),
		PreferredName:            u.PreferredName(),
		Country:                  u.Country(),
		BaseCurrency:             u.BaseCurrency(),
		Locale:                   u.Locale(),
		Timezone:                 u.Timezone(),
		UserHealth:               string(u.UserHealth()),
		UserConfidence:           string(u.UserConfidence()),
		FinancialIdentityProfile: string(u.FinancialIdentityProfile()),
		Status:                   string(u.Status()),
		UserType:                 string(u.UserType()),
		ProfileCompleteness:      u.ProfileCompleteness(),
		SourceOfTruth:            string(u.SourceOfTruth()),
		CreatedAt:                u.CreatedAt().Format(time.RFC3339),
		UpdatedAt:                u.UpdatedAt().Format(time.RFC3339),
	}
}

func toPaginatedResult(users []*domain.User, cursor string) *query.PaginatedResult {
	result := &query.PaginatedResult{
		Users:      make([]query.UserResult, 0, len(users)),
		NextCursor: cursor,
		HasMore:    cursor != "",
	}
	for _, u := range users {
		result.Users = append(result.Users, *toUserResult(u))
	}
	return result
}
