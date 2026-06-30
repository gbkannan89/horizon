package router

import (
	"github.com/gin-gonic/gin"
	"github.com/horizon/core/services/infra/api-gateway/internal/handler"
)

func RegisterRoutes(v1 *gin.RouterGroup, h *handler.Handler) {
	v1.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	// User
	ug := v1.Group("/users")
	ug.POST("", handler.CreateUser)
	ug.GET("", handler.ListUsers)
	ug.GET("/:id", handler.GetUser)
	ug.POST("/:id/activate", handler.ActivateUser)
	ug.POST("/:id/suspend", handler.SuspendUser)
	ug.POST("/:id/archive", handler.ArchiveUser)

	// Institution
	ig := v1.Group("/institutions")
	ig.POST("", handler.CreateInstitution)
	ig.GET("", handler.ListInstitutions)
	ig.GET("/:id", handler.GetInstitution)

	// Financial Event
	fg := v1.Group("/events")
	fg.POST("", handler.CreateEvent)
	fg.GET("", handler.ListEvents)
	fg.GET("/:id", handler.GetEvent)
	fg.POST("/:id/confirm", handler.ConfirmEvent)
	fg.POST("/:id/post", handler.PostEvent)
	fg.POST("/:id/reverse", handler.ReverseEvent)
	fg.POST("/:id/archive", handler.ArchiveEvent)

	// Goal
	gg := v1.Group("/goals")
	gg.POST("", handler.CreateGoal)
	gg.GET("", handler.ListGoals)
	gg.GET("/:id", handler.GetGoal)
	gg.POST("/:id/activate", handler.ActivateGoal)
	gg.POST("/:id/complete", handler.CompleteGoal)
	gg.POST("/:id/archive", handler.ArchiveGoal)

	// Account
	ag := v1.Group("/accounts")
	ag.POST("", handler.CreateAccount)
	ag.GET("", handler.ListAccounts)
	ag.GET("/:id", handler.GetAccount)
	ag.POST("/:id/activate", handler.ActivateAccount)
	ag.POST("/:id/freeze", handler.FreezeAccount)
	ag.POST("/:id/close", handler.CloseAccount)

	// Allocation
	alg := v1.Group("/allocations")
	alg.POST("", handler.CreateAllocation)
	alg.GET("", handler.ListAllocations)
	alg.GET("/:id", handler.GetAllocation)

	// Asset
	asg := v1.Group("/assets")
	asg.POST("", handler.CreateAsset)
	asg.GET("", handler.ListAssets)
	asg.GET("/:id", handler.GetAsset)

	// Liability
	lg := v1.Group("/liabilities")
	lg.POST("", handler.CreateLiability)
	lg.GET("", handler.ListLiabilities)
	lg.GET("/:id", handler.GetLiability)
	lg.POST("/:id/settle", handler.SettleLiability)

	// Portfolio
	pg := v1.Group("/portfolios")
	pg.POST("", handler.CreatePortfolio)
	pg.GET("", handler.ListPortfolios)
	pg.GET("/:id", handler.GetPortfolio)
}
