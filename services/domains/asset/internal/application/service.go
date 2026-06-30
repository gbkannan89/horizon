package application

import (
	"context"
	"fmt"
	"time"

	"github.com/horizon/core/services/domains/asset/internal/application/dto/command"
	"github.com/horizon/core/services/domains/asset/internal/application/dto/query"
	"github.com/horizon/core/services/domains/asset/internal/domain"
)

// AssetService implements all asset use cases.
type AssetService struct {
	repo domain.Repository
	pub  domain.Publisher
	fac  *domain.AssetFactory
	now  func() time.Time
}

// NewAssetService creates a new AssetService.
func NewAssetService(r domain.Repository, p domain.Publisher, n func() time.Time) *AssetService {
	return &AssetService{repo: r, pub: p, fac: domain.NewAssetFactory(), now: n}
}

func (s *AssetService) Create(ctx context.Context, cmd command.CreateAssetCommand) (*command.AssetResult, error) {
	om := cmd.OwnershipModel
	if om == "" {
		om = "Individual"
	}
	a, err := s.fac.Create("", cmd.AssetName, domain.AssetClassification(cmd.Classification),
		domain.OwnershipModel(om), cmd.OwnerID, cmd.Currency,
		domain.ValuationProfile(cmd.ValuationProfile), domain.ValuationMethod(cmd.ValuationMethod),
		cmd.OwnershipPct, cmd.SourceOfTruth, cmd.AssetType, cmd.HouseholdID,
		cmd.InstitutionID, cmd.AccountID, cmd.Quantity, cmd.ExternalIDs, cmd.Tags, cmd.Notes)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, a); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}
	if err := s.pub.Publish(domain.NewAssetCreated(a)); err != nil {
		return nil, fmt.Errorf("publish: %w", err)
	}
	return &command.AssetResult{AssetID: a.ID(), Status: string(a.Status()), Success: true}, nil
}

func (s *AssetService) Acquire(ctx context.Context, cmd command.AcquireAssetCommand) (*command.AssetResult, error) {
	a, err := s.repo.GetByID(ctx, cmd.AssetID)
	if err != nil {
		return nil, err
	}
	ad, err := time.Parse("2006-01-02", cmd.AcqDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date: %w", err)
	}
	if err := domain.AcquireAsset(a, cmd.CostBasis, ad); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, a); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}
	return &command.AssetResult{AssetID: a.ID(), Status: string(a.Status()), Success: true}, nil
}

func (s *AssetService) Activate(ctx context.Context, cmd command.ActivateAssetCommand) (*command.AssetResult, error) {
	a, err := s.repo.GetByID(ctx, cmd.AssetID)
	if err != nil {
		return nil, err
	}
	if err := domain.ActivateAsset(a); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, a); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}
	if err := s.pub.Publish(domain.NewAssetActivated(a)); err != nil {
		return nil, fmt.Errorf("publish: %w", err)
	}
	return &command.AssetResult{AssetID: a.ID(), Status: string(a.Status()), Success: true}, nil
}

func (s *AssetService) PartiallyDispose(ctx context.Context, cmd command.PartiallyDisposeCommand) (*command.AssetResult, error) {
	a, err := s.repo.GetByID(ctx, cmd.AssetID)
	if err != nil {
		return nil, err
	}
	oldQty := float64(0)
	if a.Quantity() != nil {
		oldQty = *a.Quantity()
	}
	if err := domain.PartiallyDispose(a, cmd.QtySold, cmd.Proceeds, cmd.RemainingQty); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, a); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}
	if err := s.pub.Publish(domain.NewAssetPartiallyDisposed(a, cmd.QtySold, cmd.Proceeds)); err != nil {
		return nil, fmt.Errorf("publish: %w", err)
	}
	_ = oldQty
	return &command.AssetResult{AssetID: a.ID(), Status: string(a.Status()), Success: true}, nil
}

func (s *AssetService) FullyDispose(ctx context.Context, cmd command.FullyDisposeCommand) (*command.AssetResult, error) {
	a, err := s.repo.GetByID(ctx, cmd.AssetID)
	if err != nil {
		return nil, err
	}
	dd, err := time.Parse("2006-01-02", cmd.DispDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date: %w", err)
	}
	gainLoss := cmd.Proceeds - a.CostBasis()
	if err := domain.FullyDispose(a, cmd.Proceeds, dd); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, a); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}
	if err := s.pub.Publish(domain.NewAssetFullyDisposed(a, cmd.Proceeds, gainLoss)); err != nil {
		return nil, fmt.Errorf("publish: %w", err)
	}
	return &command.AssetResult{AssetID: a.ID(), Status: string(a.Status()), Success: true}, nil
}

func (s *AssetService) Revalue(ctx context.Context, cmd command.RevalueAssetCommand) (*command.AssetResult, error) {
	a, err := s.repo.GetByID(ctx, cmd.AssetID)
	if err != nil {
		return nil, err
	}
	method := domain.ValuationMethod(cmd.Method)
	if method == "" {
		method = domain.VMMarketValue
	}
	oldVal := a.CurrentValue()
	domain.RevalueAsset(a, cmd.Value, method)
	if err := s.repo.Save(ctx, a); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}
	if err := s.pub.Publish(domain.NewAssetRevalued(a, oldVal)); err != nil {
		return nil, fmt.Errorf("publish: %w", err)
	}
	return &command.AssetResult{AssetID: a.ID(), Status: string(a.Status()), Success: true}, nil
}

func (s *AssetService) Split(ctx context.Context, cmd command.SplitAssetCommand) (*command.AssetResult, error) {
	a, err := s.repo.GetByID(ctx, cmd.AssetID)
	if err != nil {
		return nil, err
	}
	oldQty := float64(0)
	if a.Quantity() != nil {
		oldQty = *a.Quantity()
	}
	if err := domain.SplitAsset(a, cmd.OldQty, cmd.NewQty); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, a); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}
	if err := s.pub.Publish(domain.NewAssetSplit(a, oldQty, cmd.NewQty)); err != nil {
		return nil, fmt.Errorf("publish: %w", err)
	}
	return &command.AssetResult{AssetID: a.ID(), Status: string(a.Status()), Success: true}, nil
}

func (s *AssetService) Archive(ctx context.Context, cmd command.ArchiveAssetCommand) (*command.AssetResult, error) {
	a, err := s.repo.GetByID(ctx, cmd.AssetID)
	if err != nil {
		return nil, err
	}
	if err := domain.ArchiveAsset(a); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, a); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}
	if err := s.pub.Publish(domain.NewAssetArchived(a)); err != nil {
		return nil, fmt.Errorf("publish: %w", err)
	}
	return &command.AssetResult{AssetID: a.ID(), Status: string(a.Status()), Success: true}, nil
}

func toAssetView(a *domain.Asset) *query.AssetView {
	qty := a.Quantity()
	return &query.AssetView{
		AssetID: a.ID(), AssetName: a.AssetName(),
		Classification: string(a.Classification()), Currency: a.Currency(),
		Status: string(a.Status()), CostBasis: a.CostBasis(),
		CurrentValue: a.CurrentValue(), Quantity: qty, UnitPrice: a.UnitPrice(),
		OwnershipPct: a.OwnershipPercentage(), ValuationMethod: string(a.ValuationMethod()),
		LiquidityProfile: string(a.LiquidityProfile()), AssetHealth: string(a.AssetHealth()),
		CreatedAt: a.CreatedAt().Format(time.RFC3339),
	}
}

func toPag(assets []*domain.Asset, cursor string) *query.PaginatedResult {
	r := &query.PaginatedResult{Assets: make([]query.AssetView, 0, len(assets)), NextCursor: cursor, HasMore: cursor != ""}
	for _, a := range assets {
		r.Assets = append(r.Assets, *toAssetView(a))
	}
	return r
}

func (s *AssetService) GetByID(ctx context.Context, q query.GetAssetQuery) (*query.AssetView, error) {
	a, err := s.repo.GetByID(ctx, q.AssetID)
	if err != nil {
		return nil, err
	}
	return toAssetView(a), nil
}

func (s *AssetService) ListByUser(ctx context.Context, q query.ListByUserQuery) (*query.PaginatedResult, error) {
	as, c, err := s.repo.ListByUser(ctx, q.UserID, q.Cursor, q.Limit)
	if err != nil {
		return nil, err
	}
	return toPag(as, c), nil
}

func (s *AssetService) ListByClassification(ctx context.Context, q query.ListByClassQuery) (*query.PaginatedResult, error) {
	as, c, err := s.repo.ListByClassification(ctx, q.UserID, q.Classification, q.Cursor, q.Limit)
	if err != nil {
		return nil, err
	}
	return toPag(as, c), nil
}
