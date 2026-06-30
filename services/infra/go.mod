module github.com/horizon/core/services/infra

go 1.24

require (
	github.com/gomodule/redigo v1.9.2
	github.com/jackc/pgx/v5 v5.7.1
	github.com/lib/pq v1.10.9
	github.com/nats-io/nats.go v1.37.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/klauspost/compress v1.17.2 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
)

replace github.com/horizon/core/packages/errors => ../../packages/errors

replace github.com/horizon/core/packages/events => ../../packages/events
