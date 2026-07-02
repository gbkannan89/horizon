package application

import (
	"context"
	"fmt"
	"time"

	"github.com/horizon/core/services/domains/household/internal/application/dto/command"
	"github.com/horizon/core/services/domains/household/internal/application/dto/query"
	"github.com/horizon/core/services/domains/household/internal/domain"
)

type HouseholdService struct {
	repo      domain.Repository
	publisher domain.Publisher
	factory   *domain.HouseholdFactory
	now       func() time.Time
}

func NewHouseholdService(repo domain.Repository, publisher domain.Publisher, now func() time.Time) *HouseholdService {
	return &HouseholdService{
		repo:      repo,
		publisher: publisher,
		factory:   domain.NewHouseholdFactory(),
		now:       now,
	}
}

func (s *HouseholdService) Create(ctx context.Context, cmd command.CreateHouseholdCommand) (*command.HouseholdResult, error) {
	members := []domain.HouseholdMember{
		{
			UserID:       cmd.HeadOfHouseholdID,
			Role:         domain.RoleHead,
			AddedAt:      s.now().UTC().Format(time.RFC3339),
			InviteStatus: "accepted",
		},
	}

	currency := cmd.Currency
	if currency == "" {
		currency = domain.DeriveDefaultCurrency(cmd.Country)
	}

	h, err := s.factory.Create(
		"", cmd.Name,
		domain.HouseholdType(cmd.HouseholdType),
		cmd.HeadOfHouseholdID, members,
		currency, cmd.Country, cmd.Tags, cmd.Notes,
	)
	if err != nil {
		return nil, fmt.Errorf("factory: %w", err)
	}

	if err := s.repo.Save(ctx, h); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}
	if err := s.publisher.Publish(domain.NewHouseholdCreated(h)); err != nil {
		return nil, fmt.Errorf("publish: %w", err)
	}

	return &command.HouseholdResult{
		HouseholdID: h.ID(),
		Status:      string(h.Status()),
		Success:     true,
	}, nil
}

func (s *HouseholdService) Activate(ctx context.Context, cmd command.HouseholdIDCommand) (*command.HouseholdResult, error) {
	h, err := s.repo.GetByID(ctx, cmd.HouseholdID)
	if err != nil {
		return nil, err
	}
	if !h.CanTransition(cmd.RequesterID) {
		return nil, fmt.Errorf("user does not have permission to activate household")
	}
	if err := domain.ActivateHousehold(h); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, h); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}
	if err := s.publisher.Publish(domain.NewHouseholdActivated(h)); err != nil {
		return nil, fmt.Errorf("publish: %w", err)
	}
	return &command.HouseholdResult{HouseholdID: h.ID(), Status: string(h.Status()), Success: true}, nil
}

func (s *HouseholdService) Pause(ctx context.Context, cmd command.HouseholdIDCommand) (*command.HouseholdResult, error) {
	h, err := s.repo.GetByID(ctx, cmd.HouseholdID)
	if err != nil {
		return nil, err
	}
	if !h.CanTransition(cmd.RequesterID) {
		return nil, fmt.Errorf("user does not have permission to pause household")
	}
	if err := domain.PauseHousehold(h); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, h); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}
	if err := s.publisher.Publish(domain.NewHouseholdPaused(h)); err != nil {
		return nil, fmt.Errorf("publish: %w", err)
	}
	return &command.HouseholdResult{HouseholdID: h.ID(), Status: string(h.Status()), Success: true}, nil
}

func (s *HouseholdService) Resume(ctx context.Context, cmd command.HouseholdIDCommand) (*command.HouseholdResult, error) {
	h, err := s.repo.GetByID(ctx, cmd.HouseholdID)
	if err != nil {
		return nil, err
	}
	if !h.CanTransition(cmd.RequesterID) {
		return nil, fmt.Errorf("user does not have permission to resume household")
	}
	if err := domain.ResumeHousehold(h); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, h); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}
	if err := s.publisher.Publish(domain.NewHouseholdResumed(h)); err != nil {
		return nil, fmt.Errorf("publish: %w", err)
	}
	return &command.HouseholdResult{HouseholdID: h.ID(), Status: string(h.Status()), Success: true}, nil
}

func (s *HouseholdService) Dissolve(ctx context.Context, cmd command.HouseholdIDCommand) (*command.HouseholdResult, error) {
	h, err := s.repo.GetByID(ctx, cmd.HouseholdID)
	if err != nil {
		return nil, err
	}
	if !h.CanTransition(cmd.RequesterID) {
		return nil, fmt.Errorf("user does not have permission to dissolve household")
	}
	if err := domain.DissolveHousehold(h); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, h); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}
	if err := s.publisher.Publish(domain.NewHouseholdDissolved(h)); err != nil {
		return nil, fmt.Errorf("publish: %w", err)
	}
	return &command.HouseholdResult{HouseholdID: h.ID(), Status: string(h.Status()), Success: true}, nil
}

func (s *HouseholdService) Archive(ctx context.Context, cmd command.HouseholdIDCommand) (*command.HouseholdResult, error) {
	h, err := s.repo.GetByID(ctx, cmd.HouseholdID)
	if err != nil {
		return nil, err
	}
	if !h.CanTransition(cmd.RequesterID) {
		return nil, fmt.Errorf("user does not have permission to archive household")
	}
	if err := domain.ArchiveHousehold(h); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, h); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}
	if err := s.publisher.Publish(domain.NewHouseholdArchived(h)); err != nil {
		return nil, fmt.Errorf("publish: %w", err)
	}
	return &command.HouseholdResult{HouseholdID: h.ID(), Status: string(h.Status()), Success: true}, nil
}

func (s *HouseholdService) AddMember(ctx context.Context, cmd command.AddMemberCommand) (*command.HouseholdResult, error) {
	h, err := s.repo.GetByID(ctx, cmd.HouseholdID)
	if err != nil {
		return nil, err
	}
	role := domain.MemberRole(cmd.Role)
	if err := domain.ValidateMemberRole(role); err != nil {
		return nil, err
	}
	if err := h.InviteMember(cmd.RequesterID, cmd.UserID, role); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, h); err != nil {
		return nil, fmt.Errorf("add member: %w", err)
	}
	if err := s.publisher.Publish(domain.NewMemberAdded(h, cmd.UserID, role)); err != nil {
		return nil, fmt.Errorf("publish: %w", err)
	}
	return &command.HouseholdResult{HouseholdID: h.ID(), Status: string(h.Status()), Success: true}, nil
}

func (s *HouseholdService) AcceptInvite(ctx context.Context, cmd command.AcceptInviteCommand) (*command.HouseholdResult, error) {
	h, err := s.repo.GetByID(ctx, cmd.HouseholdID)
	if err != nil {
		return nil, err
	}
	if err := h.AcceptInvite(cmd.UserID); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, h); err != nil {
		return nil, fmt.Errorf("accept invite: %w", err)
	}
	return &command.HouseholdResult{HouseholdID: h.ID(), Status: string(h.Status()), Success: true}, nil
}

func (s *HouseholdService) RemoveMember(ctx context.Context, cmd command.RemoveMemberCommand) (*command.HouseholdResult, error) {
	h, err := s.repo.GetByID(ctx, cmd.HouseholdID)
	if err != nil {
		return nil, err
	}
	if err := h.RemoveMember(cmd.RequesterID, cmd.UserID); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, h); err != nil {
		return nil, fmt.Errorf("remove member: %w", err)
	}
	if err := s.publisher.Publish(domain.NewMemberRemoved(h, cmd.UserID)); err != nil {
		return nil, fmt.Errorf("publish: %w", err)
	}
	return &command.HouseholdResult{HouseholdID: h.ID(), Status: string(h.Status()), Success: true}, nil
}

func (s *HouseholdService) UpdateMemberRole(ctx context.Context, cmd command.UpdateMemberRoleCommand) (*command.HouseholdResult, error) {
	h, err := s.repo.GetByID(ctx, cmd.HouseholdID)
	if err != nil {
		return nil, err
	}
	newRole := domain.MemberRole(cmd.Role)
	if err := domain.ValidateMemberRole(newRole); err != nil {
		return nil, err
	}
	// Find old role for event
	var oldRole domain.MemberRole
	for _, m := range h.Members() {
		if m.UserID == cmd.UserID {
			oldRole = m.Role
			break
		}
	}
	if err := h.UpdateMemberRole(cmd.RequesterID, cmd.UserID, newRole); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, h); err != nil {
		return nil, fmt.Errorf("update role: %w", err)
	}
	if err := s.publisher.Publish(domain.NewMemberRoleChanged(h, cmd.UserID, oldRole, newRole)); err != nil {
		return nil, fmt.Errorf("publish: %w", err)
	}
	return &command.HouseholdResult{HouseholdID: h.ID(), Status: string(h.Status()), Success: true}, nil
}

func (s *HouseholdService) LinkAccount(ctx context.Context, cmd command.LinkAccountCommand) (*command.HouseholdResult, error) {
	h, err := s.repo.GetByID(ctx, cmd.HouseholdID)
	if err != nil { return nil, err }
	if err := h.LinkAccount(cmd.RequesterID, domain.LinkedAccount{AccountID: cmd.AccountID, AddedBy: cmd.UserID, AddedAt: s.now().UTC().Format(time.RFC3339)}); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, h); err != nil { return nil, err }
	if err := s.publisher.Publish(domain.NewAccountLinked(h, cmd.AccountID, cmd.UserID)); err != nil { return nil, err }
	return &command.HouseholdResult{HouseholdID: h.ID(), Status: string(h.Status()), Success: true}, nil
}
func (s *HouseholdService) UnlinkAccount(ctx context.Context, cmd command.UnlinkAccountCommand) (*command.HouseholdResult, error) {
	h, err := s.repo.GetByID(ctx, cmd.HouseholdID)
	if err != nil { return nil, err }
	if err := h.UnlinkAccount(cmd.RequesterID, cmd.AccountID); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, h); err != nil { return nil, err }
	if err := s.publisher.Publish(domain.NewAccountUnlinked(h, cmd.AccountID)); err != nil { return nil, err }
	return &command.HouseholdResult{HouseholdID: h.ID(), Status: string(h.Status()), Success: true}, nil
}

func (s *HouseholdService) LinkGoal(ctx context.Context, cmd command.LinkGoalCommand) (*command.HouseholdResult, error) {
	h, err := s.repo.GetByID(ctx, cmd.HouseholdID)
	if err != nil { return nil, err }
	if err := h.LinkGoal(cmd.RequesterID, domain.LinkedGoal{GoalID: cmd.GoalID, AddedBy: cmd.UserID, AddedAt: s.now().UTC().Format(time.RFC3339)}); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, h); err != nil { return nil, err }
	if err := s.publisher.Publish(domain.NewGoalLinked(h, cmd.GoalID, cmd.UserID)); err != nil { return nil, err }
	return &command.HouseholdResult{HouseholdID: h.ID(), Status: string(h.Status()), Success: true}, nil
}
func (s *HouseholdService) UnlinkGoal(ctx context.Context, cmd command.UnlinkGoalCommand) (*command.HouseholdResult, error) {
	h, err := s.repo.GetByID(ctx, cmd.HouseholdID)
	if err != nil { return nil, err }
	if err := h.UnlinkGoal(cmd.RequesterID, cmd.GoalID); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, h); err != nil { return nil, err }
	if err := s.publisher.Publish(domain.NewGoalUnlinked(h, cmd.GoalID)); err != nil { return nil, err }
	return &command.HouseholdResult{HouseholdID: h.ID(), Status: string(h.Status()), Success: true}, nil
}

func (s *HouseholdService) LinkBudget(ctx context.Context, cmd command.LinkBudgetCommand) (*command.HouseholdResult, error) {
	h, err := s.repo.GetByID(ctx, cmd.HouseholdID)
	if err != nil { return nil, err }
	if err := h.LinkBudget(cmd.RequesterID, domain.LinkedBudget{BudgetID: cmd.BudgetID, AddedBy: cmd.UserID, AddedAt: s.now().UTC().Format(time.RFC3339)}); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, h); err != nil { return nil, err }
	if err := s.publisher.Publish(domain.NewBudgetLinked(h, cmd.BudgetID, cmd.UserID)); err != nil { return nil, err }
	return &command.HouseholdResult{HouseholdID: h.ID(), Status: string(h.Status()), Success: true}, nil
}
func (s *HouseholdService) UnlinkBudget(ctx context.Context, cmd command.UnlinkBudgetCommand) (*command.HouseholdResult, error) {
	h, err := s.repo.GetByID(ctx, cmd.HouseholdID)
	if err != nil { return nil, err }
	if err := h.UnlinkBudget(cmd.RequesterID, cmd.BudgetID); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, h); err != nil { return nil, err }
	if err := s.publisher.Publish(domain.NewBudgetUnlinked(h, cmd.BudgetID)); err != nil { return nil, err }
	return &command.HouseholdResult{HouseholdID: h.ID(), Status: string(h.Status()), Success: true}, nil
}

func (s *HouseholdService) AddGoalContribution(ctx context.Context, cmd command.AddGoalContributionCommand) (*command.HouseholdResult, error) {
	h, err := s.repo.GetByID(ctx, cmd.HouseholdID)
	if err != nil { return nil, err }
	contrib := domain.GoalContribution{GoalID: cmd.GoalID, UserID: cmd.UserID, Amount: cmd.Amount, Date: s.now().UTC().Format(time.RFC3339)}
	h.AddGoalContribution(contrib)
	if err := s.repo.Save(ctx, h); err != nil { return nil, err }
	if err := s.publisher.Publish(domain.NewGoalContributionAdded(h, contrib)); err != nil { return nil, err }
	return &command.HouseholdResult{HouseholdID: h.ID(), Status: string(h.Status()), Success: true}, nil
}

func (s *HouseholdService) GetByID(ctx context.Context, q query.GetHouseholdQuery) (*query.HouseholdDetailView, error) {
	h, err := s.repo.GetByID(ctx, q.HouseholdID)
	if err != nil {
		return nil, err
	}
	return toHouseholdDetailView(h), nil
}

func (s *HouseholdService) GetHouseholdFinancialSummary(ctx context.Context, q query.GetHouseholdQuery) (*query.HouseholdFinancialSummary, error) {
	h, err := s.repo.GetByID(ctx, q.HouseholdID)
	if err != nil { return nil, err }
	
	var totalContribs int64
	for _, c := range h.GoalContributions() { totalContribs += c.Amount }

	return &query.HouseholdFinancialSummary{
		HouseholdID:        h.ID(),
		TotalAssets:        h.TotalAssets(),
		TotalLiabilities:   h.TotalLiabilities(),
		TotalNetWorth:      h.TotalNetWorth(),
		Currency:           h.Currency(),
		LinkedAccounts:     len(h.LinkedAccounts()),
		LinkedGoals:        len(h.LinkedGoals()),
		LinkedBudgets:      len(h.LinkedBudgets()),
		TotalContributions: totalContribs,
	}, nil
}

func (s *HouseholdService) ListByUser(ctx context.Context, q query.ListByUserQuery) (*query.PaginatedResult, error) {
	households, cursor, err := s.repo.ListByUser(ctx, q.UserID, q.Cursor, q.Limit)
	if err != nil {
		return nil, err
	}
	return toPaginated(households, cursor), nil
}

func (s *HouseholdService) ListByStatus(ctx context.Context, q query.ListByStatusQuery) (*query.PaginatedResult, error) {
	households, cursor, err := s.repo.ListByStatus(ctx, domain.HouseholdStatus(q.Status), q.Cursor, q.Limit)
	if err != nil {
		return nil, err
	}
	return toPaginated(households, cursor), nil
}

func toHouseholdView(h *domain.Household) query.HouseholdView {
	return query.HouseholdView{
		HouseholdID:       h.ID(),
		Name:              h.Name(),
		HouseholdType:     string(h.HouseholdType()),
		HeadOfHouseholdID: h.HeadOfHouseholdID(),
		MemberCount:       len(h.Members()),
		Status:            string(h.Status()),
		Currency:          h.Currency(),
		Country:           h.Country(),
		TotalAssets:       h.TotalAssets(),
		TotalLiabilities:  h.TotalLiabilities(),
		TotalNetWorth:     h.TotalNetWorth(),
		Health:            string(h.Health()),
		Tags:              h.Tags(),
		CreatedAt:         h.CreatedAt().Format(time.RFC3339),
		UpdatedAt:         h.UpdatedAt().Format(time.RFC3339),
	}
}

func toHouseholdDetailView(h *domain.Household) *query.HouseholdDetailView {
	v := toHouseholdView(h)
	
	members := make([]query.HouseholdMemberView, len(h.Members()))
	for i, m := range h.Members() {
		members[i] = query.HouseholdMemberView{
			UserID: m.UserID, Role: string(m.Role),
			AddedAt: m.AddedAt, InviteStatus: m.InviteStatus,
		}
	}

	linkedAccs := make([]query.LinkedAccountView, len(h.LinkedAccounts()))
	for i, a := range h.LinkedAccounts() {
		linkedAccs[i] = query.LinkedAccountView{AccountID: a.AccountID, AddedBy: a.AddedBy, AddedAt: a.AddedAt}
	}
	v.LinkedAccounts = linkedAccs

	linkedGoals := make([]query.LinkedGoalView, len(h.LinkedGoals()))
	for i, g := range h.LinkedGoals() {
		linkedGoals[i] = query.LinkedGoalView{GoalID: g.GoalID, AddedBy: g.AddedBy, AddedAt: g.AddedAt}
	}
	v.LinkedGoals = linkedGoals

	linkedBudgets := make([]query.LinkedBudgetView, len(h.LinkedBudgets()))
	for i, b := range h.LinkedBudgets() {
		linkedBudgets[i] = query.LinkedBudgetView{BudgetID: b.BudgetID, AddedBy: b.AddedBy, AddedAt: b.AddedAt}
	}
	v.LinkedBudgets = linkedBudgets

	goalContribs := make([]query.GoalContributionView, len(h.GoalContributions()))
	for i, c := range h.GoalContributions() {
		goalContribs[i] = query.GoalContributionView{GoalID: c.GoalID, UserID: c.UserID, Amount: c.Amount, Date: c.Date}
	}

	return &query.HouseholdDetailView{
		HouseholdView:     v,
		Members:           members,
		Notes:             h.Notes(),
		GoalContributions: goalContribs,
	}
}

func toPaginated(households []*domain.Household, cursor string) *query.PaginatedResult {
	r := &query.PaginatedResult{
		Households: make([]query.HouseholdView, 0, len(households)),
		NextCursor: cursor,
		HasMore:    cursor != "",
	}
	for _, h := range households {
		r.Households = append(r.Households, toHouseholdView(h))
	}
	return r
}
