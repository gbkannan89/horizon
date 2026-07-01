module github.com/horizon/core

go 1.24

require github.com/jackc/pgx/v5 v5.7.1

require (
	github.com/horizon/core/packages/http v0.0.0
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/text v0.18.0 // indirect
)

replace github.com/horizon/core/packages/errors => ./packages/errors

replace github.com/horizon/core/packages/events => ./packages/events

replace github.com/horizon/core/packages/http => ./packages/http

replace github.com/horizon/core/packages/logging => ./packages/logging

replace github.com/horizon/core/packages/config => ./packages/config

replace github.com/horizon/core/packages/auth => ./packages/auth

replace github.com/horizon/core/packages/testing => ./packages/testing

replace github.com/horizon/core/packages/telemetry => ./packages/telemetry

replace github.com/horizon/core/packages/types => ./packages/types

replace github.com/horizon/core/packages/infra => ./packages/infra
