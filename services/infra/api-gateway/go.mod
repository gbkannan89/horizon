module github.com/horizon/core/services/infra/api-gateway

go 1.24

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/horizon/core/packages/errors v0.0.0
	github.com/horizon/core/packages/events v0.0.0
	github.com/horizon/core/packages/types v0.0.0
	github.com/horizon/core/services/domains/user v0.0.0
	github.com/horizon/core/services/domains/institution v0.0.0
	github.com/horizon/core/services/domains/financial-event v0.0.0
	github.com/horizon/core/services/domains/goal v0.0.0
	github.com/horizon/core/services/domains/account v0.0.0
	github.com/horizon/core/services/domains/allocation v0.0.0
	github.com/horizon/core/services/domains/asset v0.0.0
	github.com/horizon/core/services/domains/liability v0.0.0
	github.com/horizon/core/services/domains/portfolio v0.0.0
)

replace (
	github.com/horizon/core/packages/errors => ../../../packages/errors
	github.com/horizon/core/packages/events => ../../../packages/events
	github.com/horizon/core/packages/types => ../../../packages/types
	github.com/horizon/core/services/domains/user => ../../domains/user
	github.com/horizon/core/services/domains/institution => ../../domains/institution
	github.com/horizon/core/services/domains/financial-event => ../../domains/financial-event
	github.com/horizon/core/services/domains/goal => ../../domains/goal
	github.com/horizon/core/services/domains/account => ../../domains/account
	github.com/horizon/core/services/domains/allocation => ../../domains/allocation
	github.com/horizon/core/services/domains/asset => ../../domains/asset
	github.com/horizon/core/services/domains/liability => ../../domains/liability
	github.com/horizon/core/services/domains/portfolio => ../../domains/portfolio
)
