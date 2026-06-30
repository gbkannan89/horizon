module github.com/horizon/core/services/infra/api-gateway

go 1.24

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/horizon/core/packages/errors v0.0.0
	github.com/horizon/core/packages/events v0.0.0
	github.com/horizon/core/packages/types v0.0.0
)

replace github.com/horizon/core/packages/errors => ../../../packages/errors
replace github.com/horizon/core/packages/events => ../../../packages/events
replace github.com/horizon/core/packages/types => ../../../packages/types
