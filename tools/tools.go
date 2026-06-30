//go:build tools

package tools

import (
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/vektra/mockery/v2"
	_ "github.com/pressly/goose/v3/cmd/goose"
	_ "github.com/google/wire/cmd/wire"
)
