package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/horizon/core/services/infra/api-gateway/internal/auth"
)

// Handler holds all HTTP handlers for the API.
type Handler struct {
	Auth *AuthHandler
}

// NewHandler creates a new Handler.
func NewHandler(authService *auth.Service) *Handler {
	return &Handler{Auth: NewAuthHandler(authService)}
}

func respondJSON(c *gin.Context, status int, data interface{}) { c.JSON(status, data) }
func respondCreated(c *gin.Context, data interface{})          { c.JSON(http.StatusCreated, data) }
func respondOK(c *gin.Context, data interface{})               { c.JSON(http.StatusOK, data) }
func respondNoContent(c *gin.Context)                          { c.Status(http.StatusNoContent) }

func respondError(c *gin.Context, status int, etype, code, msg string) {
	c.JSON(status, gin.H{"error": gin.H{"type": etype, "code": code, "message": msg}})
}

func respondValidation(c *gin.Context, msg string) { respondError(c, 400, "VALIDATION_ERROR", "INVALID_INPUT", msg) }
func respondNotFound(c *gin.Context, msg string)   { respondError(c, 404, "NOT_FOUND", "NOT_FOUND", msg) }
func respondConflict(c *gin.Context, msg string)   { respondError(c, 409, "CONFLICT", "CONFLICT", msg) }
func respondInternal(c *gin.Context, msg string)   { respondError(c, 500, "INTERNAL_ERROR", "INTERNAL", msg) }

// --- User Handlers ---

func (h *Handler) CreateUser(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidation(c, err.Error()); return
	}
	if svc == nil || svc.UserSvc == nil {
		respondCreated(c, gin.H{"id": "user-mock", "name": req.Name, "email": req.Email})
		return
	}
	id, err := svc.UserSvc.(UserService).Create(c.Request.Context(), req.Name, req.Email)
	if err != nil { respondInternal(c, err.Error()); return }
	respondCreated(c, gin.H{"id": id})
}

func (h *Handler) ListUsers(c *gin.Context) {
	if svc == nil || svc.UserSvc == nil {
		respondOK(c, gin.H{"data": []interface{}{}, "meta": gin.H{"total": 0}})
		return
	}
	data, err := svc.UserSvc.(UserService).List(c.Request.Context())
	if err != nil { respondInternal(c, err.Error()); return }
	respondOK(c, data)
}

func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")
	if svc == nil || svc.UserSvc == nil {
		respondOK(c, gin.H{"user_id": id, "name": "User " + id})
		return
	}
	data, err := svc.UserSvc.(UserService).Get(c.Request.Context(), id)
	if err != nil { respondNotFound(c, err.Error()); return }
	respondOK(c, data)
}

func (h *Handler) ActivateUser(c *gin.Context) {
	id := c.Param("id")
	if svc != nil && svc.UserSvc != nil {
		if err := svc.UserSvc.(UserService).Activate(c.Request.Context(), id); err != nil {
			respondConflict(c, err.Error()); return
		}
	}
	respondOK(c, gin.H{"status": "activated"})
}

func (h *Handler) SuspendUser(c *gin.Context) {
	id := c.Param("id")
	if svc != nil && svc.UserSvc != nil {
		if err := svc.UserSvc.(UserService).Suspend(c.Request.Context(), id); err != nil {
			respondConflict(c, err.Error()); return
		}
	}
	respondOK(c, gin.H{"status": "suspended"})
}

func (h *Handler) ArchiveUser(c *gin.Context) {
	id := c.Param("id")
	if svc != nil && svc.UserSvc != nil {
		if err := svc.UserSvc.(UserService).Archive(c.Request.Context(), id); err != nil {
			respondConflict(c, err.Error()); return
		}
	}
	respondOK(c, gin.H{"status": "archived"})
}

// --- Institution Handlers ---

func (h *Handler) CreateInstitution(c *gin.Context) {
	var req struct {
		Name    string `json:"name" binding:"required"`
		Type    string `json:"type" binding:"required"`
		Country string `json:"country" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil { respondValidation(c, err.Error()); return }
	if svc != nil && svc.InstitutionSvc != nil {
		id, err := svc.InstitutionSvc.(InstitutionService).Create(c.Request.Context(), req.Name, req.Type, req.Country)
		if err != nil { respondInternal(c, err.Error()); return }
		respondCreated(c, gin.H{"id": id}); return
	}
	respondCreated(c, gin.H{"id": "inst-mock", "name": req.Name, "type": req.Type, "country": req.Country})
}

func (h *Handler) ListInstitutions(c *gin.Context) {
	if svc != nil && svc.InstitutionSvc != nil {
		data, err := svc.InstitutionSvc.(InstitutionService).List(c.Request.Context())
		if err != nil { respondInternal(c, err.Error()); return }
		respondOK(c, data); return
	}
	respondOK(c, gin.H{"data": []interface{}{}, "meta": gin.H{"total": 0}})
}

func (h *Handler) GetInstitution(c *gin.Context) {
	id := c.Param("id")
	if svc != nil && svc.InstitutionSvc != nil {
		data, err := svc.InstitutionSvc.(InstitutionService).Get(c.Request.Context(), id)
		if err != nil { respondNotFound(c, err.Error()); return }
		respondOK(c, data); return
	}
	respondOK(c, gin.H{"institution_id": id, "name": "Institution " + id})
}

// --- Financial Event Handlers ---

func (h *Handler) CreateEvent(c *gin.Context) {
	var req struct {
		Type     string  `json:"type" binding:"required"`
		Amount   float64 `json:"amount" binding:"required"`
		Currency string  `json:"currency" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil { respondValidation(c, err.Error()); return }
	userID := "user-mock"
	if svc != nil && svc.FinancialEventSvc != nil {
		id, err := svc.FinancialEventSvc.(EventService).Create(c.Request.Context(), userID, req.Type, req.Amount, req.Currency)
		if err != nil { respondInternal(c, err.Error()); return }
		respondCreated(c, gin.H{"id": id}); return
	}
	respondCreated(c, gin.H{"id": "evt-mock", "type": req.Type, "amount": req.Amount, "currency": req.Currency})
}

func (h *Handler) ListEvents(c *gin.Context) {
	if svc != nil && svc.FinancialEventSvc != nil {
		data, err := svc.FinancialEventSvc.(EventService).List(c.Request.Context(), "user-mock")
		if err != nil { respondInternal(c, err.Error()); return }
		respondOK(c, data); return
	}
	respondOK(c, gin.H{"data": []interface{}{}, "meta": gin.H{"total": 0}})
}

func (h *Handler) GetEvent(c *gin.Context) {
	id := c.Param("id")
	if svc != nil && svc.FinancialEventSvc != nil {
		data, err := svc.FinancialEventSvc.(EventService).Get(c.Request.Context(), id)
		if err != nil { respondNotFound(c, err.Error()); return }
		respondOK(c, data); return
	}
	respondOK(c, gin.H{"event_id": id, "type": "mock"})
}

func (h *Handler) ConfirmEvent(c *gin.Context) {
	id := c.Param("id")
	if svc != nil && svc.FinancialEventSvc != nil {
		if err := svc.FinancialEventSvc.(EventService).Confirm(c.Request.Context(), id); err != nil {
			respondConflict(c, err.Error()); return
		}
	}
	respondOK(c, gin.H{"status": "confirmed"})
}

func (h *Handler) PostEvent(c *gin.Context) {
	id := c.Param("id")
	if svc != nil && svc.FinancialEventSvc != nil {
		if err := svc.FinancialEventSvc.(EventService).Post(c.Request.Context(), id); err != nil {
			respondConflict(c, err.Error()); return
		}
	}
	respondOK(c, gin.H{"status": "posted"})
}

func (h *Handler) ReverseEvent(c *gin.Context) {
	id := c.Param("id")
	if svc != nil && svc.FinancialEventSvc != nil {
		if err := svc.FinancialEventSvc.(EventService).Reverse(c.Request.Context(), id, "user request"); err != nil {
			respondConflict(c, err.Error()); return
		}
	}
	respondOK(c, gin.H{"status": "reversed"})
}

func (h *Handler) ArchiveEvent(c *gin.Context) {
	id := c.Param("id")
	if svc != nil && svc.FinancialEventSvc != nil {
		if err := svc.FinancialEventSvc.(EventService).Archive(c.Request.Context(), id); err != nil {
			respondConflict(c, err.Error()); return
		}
	}
	respondOK(c, gin.H{"status": "archived"})
}

// --- Goal Handlers ---

func (h *Handler) CreateGoal(c *gin.Context) {
	var req struct {
		Name         string  `json:"name" binding:"required"`
		TargetAmount float64 `json:"target_amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil { respondValidation(c, err.Error()); return }
	if svc != nil && svc.GoalSvc != nil {
		id, err := svc.GoalSvc.(GoalService).Create(c.Request.Context(), "user-mock", req.Name, req.TargetAmount)
		if err != nil { respondInternal(c, err.Error()); return }
		respondCreated(c, gin.H{"id": id}); return
	}
	respondCreated(c, gin.H{"id": "goal-mock", "name": req.Name, "target_amount": req.TargetAmount})
}

func (h *Handler) ListGoals(c *gin.Context) {
	if svc != nil && svc.GoalSvc != nil {
		data, err := svc.GoalSvc.(GoalService).List(c.Request.Context(), "user-mock")
		if err != nil { respondInternal(c, err.Error()); return }
		respondOK(c, data); return
	}
	respondOK(c, gin.H{"data": []interface{}{}, "meta": gin.H{"total": 0}})
}

func (h *Handler) GetGoal(c *gin.Context) {
	id := c.Param("id")
	if svc != nil && svc.GoalSvc != nil {
		data, err := svc.GoalSvc.(GoalService).Get(c.Request.Context(), id)
		if err != nil { respondNotFound(c, err.Error()); return }
		respondOK(c, data); return
	}
	respondOK(c, gin.H{"goal_id": id, "name": "Goal " + id, "progress": 0})
}

func (h *Handler) ActivateGoal(c *gin.Context) {
	id := c.Param("id")
	if svc != nil && svc.GoalSvc != nil {
		if err := svc.GoalSvc.(GoalService).Activate(c.Request.Context(), id); err != nil {
			respondConflict(c, err.Error()); return
		}
	}
	respondOK(c, gin.H{"status": "activated"})
}

func (h *Handler) CompleteGoal(c *gin.Context) {
	id := c.Param("id")
	if svc != nil && svc.GoalSvc != nil {
		if err := svc.GoalSvc.(GoalService).Complete(c.Request.Context(), id); err != nil {
			respondConflict(c, err.Error()); return
		}
	}
	respondOK(c, gin.H{"status": "completed"})
}

func (h *Handler) ArchiveGoal(c *gin.Context) {
	id := c.Param("id")
	if svc != nil && svc.GoalSvc != nil {
		if err := svc.GoalSvc.(GoalService).Archive(c.Request.Context(), id); err != nil {
			respondConflict(c, err.Error()); return
		}
	}
	respondOK(c, gin.H{"status": "archived"})
}

// --- Account Handlers ---

func (h *Handler) CreateAccount(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Type     string `json:"type" binding:"required"`
		Currency string `json:"currency" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil { respondValidation(c, err.Error()); return }
	if svc != nil && svc.AccountSvc != nil {
		id, err := svc.AccountSvc.(AccountService).Create(c.Request.Context(), "user-mock", req.Name, req.Type, req.Currency)
		if err != nil { respondInternal(c, err.Error()); return }
		respondCreated(c, gin.H{"id": id}); return
	}
	respondCreated(c, gin.H{"id": "acct-mock", "name": req.Name, "type": req.Type, "currency": req.Currency})
}

func (h *Handler) ListAccounts(c *gin.Context) {
	if svc != nil && svc.AccountSvc != nil {
		data, err := svc.AccountSvc.(AccountService).List(c.Request.Context(), "user-mock")
		if err != nil { respondInternal(c, err.Error()); return }
		respondOK(c, data); return
	}
	respondOK(c, gin.H{"data": []interface{}{}, "meta": gin.H{"total": 0}})
}

func (h *Handler) GetAccount(c *gin.Context) {
	id := c.Param("id")
	if svc != nil && svc.AccountSvc != nil {
		data, err := svc.AccountSvc.(AccountService).Get(c.Request.Context(), id)
		if err != nil { respondNotFound(c, err.Error()); return }
		respondOK(c, data); return
	}
	respondOK(c, gin.H{"account_id": id, "name": "Account " + id})
}

func (h *Handler) ActivateAccount(c *gin.Context) {
	id := c.Param("id")
	if svc != nil && svc.AccountSvc != nil {
		if err := svc.AccountSvc.(AccountService).Activate(c.Request.Context(), id); err != nil {
			respondConflict(c, err.Error()); return
		}
	}
	respondOK(c, gin.H{"status": "activated"})
}

func (h *Handler) FreezeAccount(c *gin.Context) {
	id := c.Param("id")
	if svc != nil && svc.AccountSvc != nil {
		if err := svc.AccountSvc.(AccountService).Freeze(c.Request.Context(), id); err != nil {
			respondConflict(c, err.Error()); return
		}
	}
	respondOK(c, gin.H{"status": "frozen"})
}

func (h *Handler) CloseAccount(c *gin.Context) {
	id := c.Param("id")
	if svc != nil && svc.AccountSvc != nil {
		if err := svc.AccountSvc.(AccountService).Close(c.Request.Context(), id); err != nil {
			respondConflict(c, err.Error()); return
		}
	}
	respondOK(c, gin.H{"status": "closed"})
}

// --- Allocation Handlers ---

func (h *Handler) CreateAllocation(c *gin.Context) {
	var req struct {
		GoalID   string  `json:"goal_id" binding:"required"`
		SourceID string  `json:"source_id" binding:"required"`
		Type     string  `json:"type" binding:"required"`
		Currency string  `json:"currency" binding:"required"`
		Amount   float64 `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil { respondValidation(c, err.Error()); return }
	if svc != nil && svc.AllocationSvc != nil {
		id, err := svc.AllocationSvc.(AllocationService).Create(c.Request.Context(), req.GoalID, req.SourceID, req.Type, req.Currency, req.Amount)
		if err != nil { respondInternal(c, err.Error()); return }
		respondCreated(c, gin.H{"id": id}); return
	}
	respondCreated(c, gin.H{"id": "alloc-mock", "goal_id": req.GoalID, "amount": req.Amount})
}

func (h *Handler) ListAllocations(c *gin.Context) {
	if svc != nil && svc.AllocationSvc != nil {
		data, err := svc.AllocationSvc.(AllocationService).List(c.Request.Context())
		if err != nil { respondInternal(c, err.Error()); return }
		respondOK(c, data); return
	}
	respondOK(c, gin.H{"data": []interface{}{}, "meta": gin.H{"total": 0}})
}

func (h *Handler) GetAllocation(c *gin.Context) {
	id := c.Param("id")
	if svc != nil && svc.AllocationSvc != nil {
		data, err := svc.AllocationSvc.(AllocationService).Get(c.Request.Context(), id)
		if err != nil { respondNotFound(c, err.Error()); return }
		respondOK(c, data); return
	}
	respondOK(c, gin.H{"allocation_id": id})
}

// --- Asset Handlers ---

func (h *Handler) CreateAsset(c *gin.Context) {
	var req struct {
		Name     string  `json:"name" binding:"required"`
		Class    string  `json:"classification" binding:"required"`
		Currency string  `json:"currency" binding:"required"`
		Value    float64 `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil { respondValidation(c, err.Error()); return }
	if svc != nil && svc.AssetSvc != nil {
		id, err := svc.AssetSvc.(AssetService).Create(c.Request.Context(), "user-mock", req.Name, req.Class, req.Currency, req.Value)
		if err != nil { respondInternal(c, err.Error()); return }
		respondCreated(c, gin.H{"id": id}); return
	}
	respondCreated(c, gin.H{"id": "asset-mock", "name": req.Name, "value": req.Value})
}

func (h *Handler) ListAssets(c *gin.Context) {
	if svc != nil && svc.AssetSvc != nil {
		data, err := svc.AssetSvc.(AssetService).List(c.Request.Context(), "user-mock")
		if err != nil { respondInternal(c, err.Error()); return }
		respondOK(c, data); return
	}
	respondOK(c, gin.H{"data": []interface{}{}, "meta": gin.H{"total": 0}})
}

func (h *Handler) GetAsset(c *gin.Context) {
	id := c.Param("id")
	if svc != nil && svc.AssetSvc != nil {
		data, err := svc.AssetSvc.(AssetService).Get(c.Request.Context(), id)
		if err != nil { respondNotFound(c, err.Error()); return }
		respondOK(c, data); return
	}
	respondOK(c, gin.H{"asset_id": id, "name": "Asset " + id})
}

// --- Liability Handlers ---

func (h *Handler) CreateLiability(c *gin.Context) {
	var req struct {
		Name      string  `json:"name" binding:"required"`
		Class     string  `json:"classification" binding:"required"`
		Currency  string  `json:"currency" binding:"required"`
		Principal float64 `json:"principal" binding:"required"`
		Rate      float64 `json:"interest_rate" binding:"required"`
		Maturity  string  `json:"maturity_date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil { respondValidation(c, err.Error()); return }
	if svc != nil && svc.LiabilitySvc != nil {
		mt, _ := time.Parse("2006-01-02", req.Maturity)
		id, err := svc.LiabilitySvc.(LiabilityService).Create(c.Request.Context(), "user-mock", req.Name, req.Class, req.Currency, req.Principal, req.Rate, mt)
		if err != nil { respondInternal(c, err.Error()); return }
		respondCreated(c, gin.H{"id": id}); return
	}
	respondCreated(c, gin.H{"id": "liab-mock", "name": req.Name, "principal": req.Principal})
}

func (h *Handler) ListLiabilities(c *gin.Context) {
	if svc != nil && svc.LiabilitySvc != nil {
		data, err := svc.LiabilitySvc.(LiabilityService).List(c.Request.Context(), "user-mock")
		if err != nil { respondInternal(c, err.Error()); return }
		respondOK(c, data); return
	}
	respondOK(c, gin.H{"data": []interface{}{}, "meta": gin.H{"total": 0}})
}

func (h *Handler) GetLiability(c *gin.Context) {
	id := c.Param("id")
	if svc != nil && svc.LiabilitySvc != nil {
		data, err := svc.LiabilitySvc.(LiabilityService).Get(c.Request.Context(), id)
		if err != nil { respondNotFound(c, err.Error()); return }
		respondOK(c, data); return
	}
	respondOK(c, gin.H{"liability_id": id, "name": "Liability " + id})
}

func (h *Handler) SettleLiability(c *gin.Context) {
	id := c.Param("id")
	if svc != nil && svc.LiabilitySvc != nil {
		if err := svc.LiabilitySvc.(LiabilityService).Settle(c.Request.Context(), id); err != nil {
			respondConflict(c, err.Error()); return
		}
	}
	respondOK(c, gin.H{"status": "settled"})
}

// --- Portfolio Handlers ---

func (h *Handler) CreatePortfolio(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Type     string `json:"type" binding:"required"`
		Currency string `json:"currency" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil { respondValidation(c, err.Error()); return }
	if svc != nil && svc.PortfolioSvc != nil {
		id, err := svc.PortfolioSvc.(PortfolioService).Create(c.Request.Context(), "user-mock", req.Name, req.Type, req.Currency)
		if err != nil { respondInternal(c, err.Error()); return }
		respondCreated(c, gin.H{"id": id}); return
	}
	respondCreated(c, gin.H{"id": "pf-mock", "name": req.Name, "type": req.Type})
}

func (h *Handler) ListPortfolios(c *gin.Context) {
	if svc != nil && svc.PortfolioSvc != nil {
		data, err := svc.PortfolioSvc.(PortfolioService).List(c.Request.Context(), "user-mock")
		if err != nil { respondInternal(c, err.Error()); return }
		respondOK(c, data); return
	}
	respondOK(c, gin.H{"data": []interface{}{}, "meta": gin.H{"total": 0}})
}

func (h *Handler) GetPortfolio(c *gin.Context) {
	id := c.Param("id")
	if svc != nil && svc.PortfolioSvc != nil {
		data, err := svc.PortfolioSvc.(PortfolioService).Get(c.Request.Context(), id)
		if err != nil { respondNotFound(c, err.Error()); return }
		respondOK(c, data); return
	}
	respondOK(c, gin.H{"portfolio_id": id, "name": "Portfolio " + id})
}
