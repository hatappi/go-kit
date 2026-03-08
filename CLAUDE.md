# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test Commands

- **Run all tests:** `make test` (runs `go test ./... -v`)
- **Run a single test:** `go test -v -run TestFuncName ./path/to/package/`
- **Tidy dependencies:** `go mod tidy`

## Architecture

This is a Go utility library (`github.com/hatappi/go-kit`) providing reusable packages. Go 1.22+, toolchain 1.23.3.

### Packages

- **storage/** — File storage abstraction with a `Storage` interface (`Save`, `Get`, `Delete`, `Ping`). Factory function `NewStorage()` selects provider based on config type. Providers: `provider/disk.go` (local filesystem), `provider/s3.go` (AWS S3 via SDK v2). Functional options pattern for `Save()` via `option/` package.
- **log/** — Context-aware logging using `go-logr/logr` as the facade. `log/zap/` provides a zap-backed implementation. `log/retryablehttp/` provides an slog-based logger adapter for `go-retryablehttp`. Logger is stored/retrieved from `context.Context` via `WithContext`/`FromContext`.
- **line/notify/** — LINE Notify API client for sending text notifications.

### Key Patterns

- **Interface-based design** — Storage providers and S3 client are behind interfaces for testability (mock S3 in tests).
- **Context-first** — All I/O operations take `context.Context` as the first parameter.
- **Functional options** — `SaveOptionFunc` for extensible method configuration without breaking changes.
- **Config via env tags** — Storage config uses `envconfig` struct tags for environment variable binding.
- **Table-driven tests** with `google/go-cmp` for comparisons.
