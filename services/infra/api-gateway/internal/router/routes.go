package router

import (
	"github.com/gin-gonic/gin"
	"github.com/horizon/core/services/infra/api-gateway/internal/handler"
	"github.com/horizon/core/services/infra/api-gateway/internal/middleware"
)

func RegisterRoutes(v1 *gin.RouterGroup, h *handler.Handler) {
	// Public routes
	v1.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	// Auth routes (public)
	auth := v1.Group("/auth")
	auth.POST("/login", h.Auth.Login)
	auth.POST("/refresh", h.Auth.Refresh)
	auth.POST("/logout", h.Auth.Logout)
	auth.GET("/me", h.Auth.Me)

	// Protected routes — require authentication
	prot := v1.Group("")
	prot.Use(middleware.Authentication) // validates JWT from Authorization header

	// User routes
	ug := prot.Group("/users")
	ug.POST("", h.CreateUser)
	ug.GET("", h.ListUsers)
	ug.GET("/:id", h.GetUser)
	ug.POST("/:id/activate", h.ActivateUser)
	ug.POST("/:id/suspend", h.SuspendUser)
	ug.POST("/:id/archive", h.ArchiveUser)

	// Institution routes
	ig := prot.Group("/institutions")
	ig.POST("", h.CreateInstitution)
	ig.GET("", h.ListInstitutions)
	ig.GET("/:id", h.GetInstitution)

	// Financial Event routes
	fg := prot.Group("/events")
	fg.POST("", h.CreateEvent)
	fg.GET("", h.ListEvents)
	fg.GET("/:id", h.GetEvent)
	fg.POST("/:id/confirm", h.ConfirmEvent)
	fg.POST("/:id/post", h.PostEvent)
	fg.POST("/:id/reverse", h.ReverseEvent)
	fg.POST("/:id/archive", h.ArchiveEvent)

	// Goal routes
	gg := prot.Group("/goals")
	gg.POST("", h.CreateGoal)
	gg.GET("", h.ListGoals)
	gg.GET("/:id", h.GetGoal)
	gg.POST("/:id/activate", h.ActivateGoal)
	gg.POST("/:id/complete", h.CompleteGoal)
	gg.POST("/:id/archive", h.ArchiveGoal)

	// Account routes
	ag := prot.Group("/accounts")
	ag.POST("", h.CreateAccount)
	ag.GET("", h.ListAccounts)
	ag.GET("/:id", h.GetAccount)
	ag.POST("/:id/activate", h.ActivateAccount)
	ag.POST("/:id/freeze", h.FreezeAccount)
	ag.POST("/:id/close", h.CloseAccount)

	// Allocation routes
	alg := prot.Group("/allocations")
	alg.POST("", h.CreateAllocation)
	alg.GET("", h.ListAllocations)
	alg.GET("/:id", h.GetAllocation)

	// Asset routes
	asg := prot.Group("/assets")
	asg.POST("", h.CreateAsset)
	asg.GET("", h.ListAssets)
	asg.GET("/:id", h.GetAsset)

	// Liability routes
	lg := prot.Group("/liabilities")
	lg.POST("", h.CreateLiability)
	lg.GET("", h.ListLiabilities)
	lg.GET("/:id", h.GetLiability)
	lg.POST("/:id/settle", h.SettleLiability)

	// Portfolio routes
	pg := prot.Group("/portfolios")
	pg.POST("", h.CreatePortfolio)
	pg.GET("", h.ListPortfolios)
	pg.GET("/:id", h.GetPortfolio)
}
