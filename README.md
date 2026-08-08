# MCP VictoriaMetrics Demo

This project is a small Go-based MCP setup with:

- an MCP **server** in `cmd/server` exposing tools backed by Prometheus/VictoriaMetrics API
- an MCP **client** in `cmd/client` that connects to the server and calls tools

## Tools exposed by server

- `buildInfo` - returns VictoriaMetrics build information
- `series` - calls `vmApi.Series` to fetch matching label sets
- `metadata` - calls `vmApi.Metadata` to fetch metric metadata

## Project layout

- `cmd/server/main.go` - MCP server and tool/prompt registration
- `cmd/client/main.go` - MCP client example calls
- `Makefile` - convenience commands

## Commands

Run all commands from the project root (`VM`).

### Build server binary

```bash
make build-server
```

This builds `cmd/server/server`.

### Run client

```bash
make run-client
```

Client authenticates using `MCP_API_KEY` (default: `local-dev-mcp-key`).

### Build and run server in Docker

Build image locally:

```bash
docker compose build mcp-server
```

Run the MCP server container:

```bash
docker compose up --build mcp-server
```

By default, containerized server connects to VictoriaMetrics at `http://host.docker.internal:8428`.
Override it with:

```bash
VM_ADDRESS=http://localhost:8428 docker compose up --build mcp-server
```

The MCP Streamable HTTP endpoint is exposed at `http://localhost:8080/mcp`.
Calls to this endpoint require the `X-API-Key` header.
Prometheus metrics are exposed at `http://localhost:8080/metrics`.

## Notes
- Server uses `VM_ADDRESS` when set; otherwise defaults to `http://localhost:8428`.
- Server uses `MCP_HTTP_ADDR` when set; otherwise defaults to `:8080`.
- Server and client use `MCP_API_KEY` for API key authentication (default: `local-dev-mcp-key`).
- TLS verification is currently disabled in `cmd/server/main.go` for local/dev usage.
- Server exports MCP metrics for tool calls, tool errors, tool latency, HTTP requests, HTTP errors, and HTTP latency.

