module github.com/horizon/core/tests/e2e

go 1.25.2

require github.com/horizon/core/services/domains/household v0.0.0

require (
	github.com/horizon/core v0.0.0-00010101000000-000000000000 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.1 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/text v0.19.0 // indirect
)

replace github.com/horizon/core => ../../
replace github.com/horizon/core/packages/errors => ../../packages/errors
replace github.com/horizon/core/packages/events => ../../packages/events
replace github.com/horizon/core/services/domains/household => ../../services/domains/household
