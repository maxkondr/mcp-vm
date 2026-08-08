# AGENTS Guide for `mcp-vm`

## Big picture
- This repo is a small MCP-over-HTTP service that proxies selected VictoriaMetrics/Prometheus API calls as MCP tools.
- Runtime boundary:
  - Server entrypoint: `cmd/server/main.go`
  - MCP composition: `server/server.go`
  - Tool handlers + VM API adapter: `server/mcp/*.go`
  - Demo client: `cmd/client/main.go`
- Request flow: HTTP `/mcp` -> API key middleware (`server/auth/auth.go`) -> MCP handler -> tool function (`server/mcp/tools.go`) -> package-global `vmApi` client (`server/mcp/vm.go`) -> VictoriaMetrics.

## Critical implementation patterns
- MCP server assembly is centralized in `server.NewMCPServer()`; add new tools/prompts through `server/mcp/registration.go`.
- Tool handlers are typed `mcp.AddTool` functions returning `(result, structuredOutput, error)`; this repo uses structured output structs (see `BuildInfoOutput`, `SeriesOutput`, `MetadataOutput` in `server/mcp/tools.go`).
- Every tool currently applies explicit metrics + logs:
  - start: `metrics.ObserveToolStart("<tool>")`
  - end: `metrics.ObserveToolEnd("<tool>", start, err)`
  - logrus `Info/Error` with structured fields
- HTTP-level metrics are middleware-based (`metrics.WithHTTPMetrics`) and mounted around the mux in `cmd/server/main.go`.
- Authentication convention is fixed to header `X-API-Key` (`server/config/config.go`), with constant-time compare in middleware.

## Non-obvious gotchas
- `server/mcp/vm.go` holds a package-global `vmApi`; tool handlers assume it is initialized before first call.
- `InitVictoriaMetricsApi(...)` exists in `server/mcp/vm.go`; when changing startup/wiring, ensure this init path is preserved or tools will fail/panic on nil API usage.
- `cmd/server/main.go` intentionally sets `InsecureSkipVerify: true` for default HTTP client (local/dev posture).
- Build flags in `Makefile` and `Dockerfile` include `-gcflags="all=-N -l"` (disables optimizations/inlining), so binaries are debug-friendly but not perf-optimized.

## Developer workflows used in this repo
- Build binaries (output to `bin/`):
  - `make build-server`
  - `make build-client`
- Run baseline package checks:
  - `go test ./...` (currently no test files; use as compile/regression gate)
- Containerized server path:
  - `docker compose build mcp-server`
  - `docker compose up --build mcp-server`
- Local defaults and env contract:
  - `VM_ADDRESS` default `http://localhost:8428`
  - `MCP_HTTP_ADDR` default `:8080`
  - `MCP_API_KEY` default `local-dev-mcp-key`
  - Client endpoint env: `MCP_SERVER_URL` default `http://localhost:8080/mcp`

## Integration points
- MCP SDK: `github.com/modelcontextprotocol/go-sdk/mcp` for server/client transport and tool registration.
- VictoriaMetrics/Prometheus API client: `github.com/prometheus/client_golang/api/prometheus/v1` (`Buildinfo`, `Series`, `Metadata`).
- Metrics export: Prometheus `/metrics` endpoint via `promhttp` in `cmd/server/main.go`.
- If adding a new VM-backed tool, mirror existing pattern in `server/mcp/tools.go` + register it in `server/mcp/registration.go` + preserve metrics/log/auth behavior.

