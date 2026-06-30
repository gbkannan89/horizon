module github.com/horizon/core/services/engines/recommendation

go 1.24

require github.com/jackc/pgx/v5 v5.7.1

require github.com/horizon/core/packages/errors v0.0.0
require github.com/horizon/core/packages/events v0.0.0

replace github.com/horizon/core/packages/errors => ../../../packages/errors
replace github.com/horizon/core/packages/events => ../../../packages/events

