module github.com/horizon/core/tools

go 1.24

tool (
	github.com/bufbuild/buf/cmd/buf
	github.com/golangci/golangci-lint/cmd/golangci-lint
	github.com/vektra/mockery/v2
	github.com/pressly/goose/v3/cmd/goose
	github.com/google/wire/cmd/wire
)
