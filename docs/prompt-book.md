# Prompt Book

This document describes system-level prompt behavior and guardrails for the MCP server in this repository.

## 1) System prompts exposed by MCP server

Prompt templates are registered in `server/mcp/registration.go` and implemented in `server/mcp/prompts.go`.

### `vm_buildinfo`

- Description: prompt template for `buildInfo` tool.
- User message template:
  - "Call the buildInfo tool and provide the result as a short, readable summary."

### `vm_series`

- Description: prompt template for `series` tool.
- Arguments:
  - `match` (required): PromQL selector
  - `start` (optional): RFC3339 timestamp
  - `end` (optional): RFC3339 timestamp
- User message behavior:
  - Always asks to call `series` with `match`.
  - Adds `start`/`end` only when provided.

### `vm_metadata`

- Description: prompt template for `metadata` tool.
- Arguments:
  - `metric` (optional): metric name filter
  - `limit` (optional): max metadata entries per metric
- User message behavior:
  - Calls `metadata` with no arguments by default.
  - Includes `metric`/`limit` when provided.

## 2) Guardrails and behavioral constraints

## Access control

- All `/mcp` requests must include header `X-API-Key`.
- Rejected requests return HTTP `401 Unauthorized`.
- API key comparison uses constant-time compare in `server/auth/auth.go`.

## Tool execution boundaries

- Only explicitly registered tools can be called:
  - `buildInfo`
  - `series`
  - `metadata`
- Tools are read-only over VictoriaMetrics API and do not mutate external state.

## Observability constraints

- Every tool records start/end metrics and emits structured logs.
- HTTP-level metrics are exported through middleware to `/metrics`.

## Runtime safety notes

- VictoriaMetrics client is initialized at startup via `InitVictoriaMetricsApi(...)`.
- Local/dev mode currently disables TLS verification in the default HTTP client.

## 3) Intended agent behavior with this MCP server

- Use prompt templates to generate consistent tool-call instructions.
- Keep responses concise and grounded in tool output.
- Do not infer metrics not returned by VictoriaMetrics APIs.
- On tool errors, surface a clear failure message and suggest retry with corrected arguments.

