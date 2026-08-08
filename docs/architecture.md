# MCP Server Architecture

```mermaid
flowchart LR
    U[LLM Agent / MCP Client] -->|HTTP /mcp + X-API-Key| A[API Key Middleware]
    A --> H[MCP Streamable HTTP Handler]
    H --> T[Tool Router]

    T --> BI[buildInfo tool]
    T --> SE[series tool]
    T --> ME[metadata tool]

    BI --> V[(VictoriaMetrics / Prometheus API)]
    SE --> V
    ME --> V

    H --> P[Prompt Templates]
    P --> PB[vm_buildinfo / vm_series / vm_metadata]

    T --> M[Tool Metrics + Logs]
    A --> HM[HTTP Metrics + Auth Logs]
    HM --> PR[/metrics Prometheus Exporter/]

    subgraph Agent Runtime Context
        RAG["RAG Layer (optional external)"]
        MEM["Conversation Memory (optional external)"]
    end

    U -. context retrieval .-> RAG
    U -. session memory .-> MEM
```

## Interaction model

- MCP server is a tool gateway: it does not implement vector storage or memory persistence internally.
- RAG and memory layers are external to this service and are expected to be managed by the MCP client or orchestrating agent runtime.
- The server enforces API-key authentication, then executes only registered tools against VictoriaMetrics.
- Prometheus-compatible runtime metrics are exposed on `/metrics`.

