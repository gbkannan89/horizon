module github.com/horizon/core/services/domains/user

go 1.24

require github.com/horizon/core/packages/errors v0.0.0
require github.com/horizon/core/packages/events v0.0.0

replace github.com/horizon/core/packages/errors => ../../../packages/errors
replace github.com/horizon/core/packages/events => ../../../packages/events
