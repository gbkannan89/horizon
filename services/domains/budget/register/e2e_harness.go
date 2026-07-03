package register

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/horizon/core/services/domains/budget/internal/application"
	"github.com/horizon/core/services/domains/budget/internal/domain"
)

// E2ETestHarness provides an in-memory test harness for E2E tests.
// All methods return only standard Go types so external test packages
// can use them without importing internal packages.
type E2ETestHarness struct {
	repo *e2eMockRepo
}

// NewE2ETestHarness creates a new test harness with an empty in-memory repo.
func NewE2ETestHarness() *E2ETestHarness {
	return &E2ETestHarness{repo: newE2eMockRepo()}
}

// RegisterRoutesOnMux registers all budget HTTP routes on the given mux using
// the harness's in-memory repository and a no-op event publisher.
func (h *E2ETestHarness) RegisterRoutesOnMux(mux *http.ServeMux) {
	pub := &noopPublisher{}
	svc := application.NewBudgetService(h.repo, pub, time.Now)
	handler := &budgetHandler{svc: svc}

	mux.HandleFunc("POST /api/v1/budgets", handler.createBudget)
	mux.HandleFunc("GET /api/v1/budgets/{id}", handler.getBudget)
	mux.HandleFunc("GET /api/v1/budgets", handler.listBudgets)
	mux.HandleFunc("POST /api/v1/budgets/{id}/activate", handler.activateBudget)
	mux.HandleFunc("POST /api/v1/budgets/{id}/pause", handler.pauseBudget)
	mux.HandleFunc("POST /api/v1/budgets/{id}/resume", handler.resumeBudget)
	mux.HandleFunc("POST /api/v1/budgets/{id}/complete", handler.completeBudget)
	mux.HandleFunc("POST /api/v1/budgets/{id}/archive", handler.archiveBudget)
	mux.HandleFunc("PUT /api/v1/budgets/{id}/category", handler.updateCategory)
	mux.HandleFunc("GET /api/v1/budgets/{id}/vs-actual", handler.budgetVsActual)
	mux.HandleFunc("GET /api/v1/households/{id}/budgets", handler.listByHousehold)
}

// BudgetExists returns true if a budget with the given ID is stored.
func (h *E2ETestHarness) BudgetExists(id string) bool {
	h.repo.mu.Lock()
	defer h.repo.mu.Unlock()
	_, ok := h.repo.budgets[id]
	return ok
}

// GetBudgetStatus returns the status string of the budget, or empty string if not found.
func (h *E2ETestHarness) GetBudgetStatus(id string) string {
	h.repo.mu.Lock()
	defer h.repo.mu.Unlock()
	b, ok := h.repo.budgets[id]
	if !ok {
		return ""
	}
	return string(b.Status())
}

// GetCategoryID returns the ID of the first category matching categoryName in the budget.
func (h *E2ETestHarness) GetCategoryID(budgetID, categoryName string) string {
	h.repo.mu.Lock()
	defer h.repo.mu.Unlock()
	b, ok := h.repo.budgets[budgetID]
	if !ok {
		return ""
	}
	for _, c := range b.Categories() {
		if c.Category == categoryName {
			return c.ID
		}
	}
	return ""
}

// GetCategorySpent returns the spent amount for a named category, or 0.
func (h *E2ETestHarness) GetCategorySpent(budgetID, categoryName string) int64 {
	h.repo.mu.Lock()
	defer h.repo.mu.Unlock()
	b, ok := h.repo.budgets[budgetID]
	if !ok {
		return 0
	}
	for _, c := range b.Categories() {
		if c.Category == categoryName {
			return c.SpentAmount
		}
	}
	return 0
}

// GetCategoryRemaining returns the remaining amount for a named category, or 0.
func (h *E2ETestHarness) GetCategoryRemaining(budgetID, categoryName string) int64 {
	h.repo.mu.Lock()
	defer h.repo.mu.Unlock()
	b, ok := h.repo.budgets[budgetID]
	if !ok {
		return 0
	}
	for _, c := range b.Categories() {
		if c.Category == categoryName {
			return c.RemainingAmount
		}
	}
	return 0
}

// GetTotalSpent returns the total_spent from the budget entity.
func (h *E2ETestHarness) GetTotalSpent(budgetID string) int64 {
	h.repo.mu.Lock()
	defer h.repo.mu.Unlock()
	b, ok := h.repo.budgets[budgetID]
	if !ok {
		return 0
	}
	return b.TotalSpent()
}

// GetTotalBudgeted returns the total_budgeted from the budget entity.
func (h *E2ETestHarness) GetTotalBudgeted(budgetID string) int64 {
	h.repo.mu.Lock()
	defer h.repo.mu.Unlock()
	b, ok := h.repo.budgets[budgetID]
	if !ok {
		return 0
	}
	return b.TotalBudgeted()
}

// ComputeTotalSpentFromCategories sums spent amounts from all categories.
func (h *E2ETestHarness) ComputeTotalSpentFromCategories(budgetID string) int64 {
	h.repo.mu.Lock()
	defer h.repo.mu.Unlock()
	b, ok := h.repo.budgets[budgetID]
	if !ok {
		return 0
	}
	var total int64
	for _, c := range b.Categories() {
		total += c.SpentAmount
	}
	return total
}

// ComputeTotalBudgetedFromCategories sums budgeted amounts from all categories.
func (h *E2ETestHarness) ComputeTotalBudgetedFromCategories(budgetID string) int64 {
	h.repo.mu.Lock()
	defer h.repo.mu.Unlock()
	b, ok := h.repo.budgets[budgetID]
	if !ok {
		return 0
	}
	var total int64
	for _, c := range b.Categories() {
		total += c.BudgetedAmount
	}
	return total
}

// ---------------------------------------------------------------------------
// In-memory mock repository (internal to the harness)
// ---------------------------------------------------------------------------

type e2eMockRepo struct {
	mu      sync.Mutex
	budgets map[string]*domain.Budget
}

func newE2eMockRepo() *e2eMockRepo {
	return &e2eMockRepo{budgets: make(map[string]*domain.Budget)}
}

func (r *e2eMockRepo) Save(_ context.Context, b *domain.Budget) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := b.ID()
	if id == "" {
		id = fmt.Sprintf("budget-%d", len(r.budgets)+1)
		now := time.Now().UTC()

		// Assign IDs to categories that lack them
		cats := b.Categories()
		for i := range cats {
			if cats[i].ID == "" {
				cats[i].ID = fmt.Sprintf("%s-cat-%d", id, i+1)
			}
		}

		rebuilt := domain.ReconstructFromDB(
			id, b.UserID(), b.HouseholdID(), b.Name(), b.Period(),
			b.StartDate(), b.EndDate(), b.Status(),
			b.TotalBudgeted(), b.TotalSpent(), b.TotalRemaining(),
			b.Currency(), cats, b.Tags(), now, now,
		)
		*b = *rebuilt
	}
	r.budgets[b.ID()] = b
	return nil
}

func (r *e2eMockRepo) UpdateStatus(_ context.Context, id string, from, to domain.BudgetStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.budgets[id]
	if !ok {
		return fmt.Errorf("budget not found")
	}
	if b.Status() != from {
		return fmt.Errorf("status mismatch")
	}
	b.SetStatus(to)
	return nil
}

func (r *e2eMockRepo) GetByID(_ context.Context, id string) (*domain.Budget, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.budgets[id]
	if !ok {
		return nil, fmt.Errorf("budget not found")
	}
	return b, nil
}

func (r *e2eMockRepo) ListByUser(_ context.Context, userID string, _ string, _ int) ([]*domain.Budget, string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.Budget
	for _, b := range r.budgets {
		if b.UserID() == userID {
			result = append(result, b)
		}
	}
	return result, "", nil
}

func (r *e2eMockRepo) ListByPeriod(_ context.Context, _ string, _, _ time.Time, _ string, _ int) ([]*domain.Budget, string, error) {
	return nil, "", nil
}

func (r *e2eMockRepo) GetByCategory(_ context.Context, _, _ string, _ string, _ int) ([]*domain.Budget, string, error) {
	return nil, "", nil
}

func (r *e2eMockRepo) GetBudgetVsActual(_ context.Context, budgetID string) ([]domain.BudgetCategory, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.budgets[budgetID]
	if !ok {
		return nil, fmt.Errorf("budget not found")
	}
	return b.Categories(), nil
}

func (r *e2eMockRepo) ListByHousehold(_ context.Context, householdID string, _ string, _ int) ([]*domain.Budget, string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.Budget
	for _, b := range r.budgets {
		if b.HouseholdID() == householdID {
			result = append(result, b)
		}
	}
	return result, "", nil
}

// Compile-time check
var _ domain.Repository = (*e2eMockRepo)(nil)
