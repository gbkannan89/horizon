package engine

import (
	"context"

	"github.com/horizon/core/services/experiences/data/internal/api"
	"github.com/horizon/core/services/experiences/data/internal/infrastructure/persistence"
)

type Exporter struct {
	provider *persistence.DataProvider
}

func NewExporter(provider *persistence.DataProvider) *Exporter {
	return &Exporter{provider: provider}
}

func (e *Exporter) ExportAll(ctx context.Context, userID string) api.ExportData {
	data := api.NewExportData(userID)
	data.Accounts = e.provider.GetAccounts(ctx, userID)
	data.Transactions = e.provider.GetTransactions(ctx, userID)
	data.Goals = e.provider.GetGoals(ctx, userID)
	data.Allocations = e.provider.GetAllocations(ctx, userID)
	data.Assets = e.provider.GetAssets(ctx, userID)
	data.Liabilities = e.provider.GetLiabilities(ctx, userID)
	data.Portfolios = e.provider.GetPortfolios(ctx, userID)
	data.Institutions = e.provider.GetInstitutions(ctx, userID)
	data.Households = e.provider.GetHouseholds(ctx, userID)
	data.HealthScores = e.provider.GetHealthScores(ctx, userID)
	data.RiskAssessments = e.provider.GetRiskAssessments(ctx, userID)
	data.Recommendations = e.provider.GetRecommendations(ctx, userID)
	data.Optimizations = e.provider.GetOptimizations(ctx, userID)
	data.Projections = e.provider.GetProjections(ctx, userID)
	data.Simulations = e.provider.GetSimulations(ctx, userID)
	return data
}
