package config

import "os"

// APIKeyHeader is the HTTP header name used for MCP API key authentication.
const APIKeyHeader = "X-API-Key"

// VMAddressFromEnv returns the VictoriaMetrics base URL from VM_ADDRESS or a local default.
func VMAddressFromEnv() string {
    if v := os.Getenv("VM_ADDRESS"); v != "" {
        return v
    }
    return "http://localhost:8428"
}

// HTTPAddressFromEnv returns the HTTP listen address from MCP_HTTP_ADDR or a local default.
func HTTPAddressFromEnv() string {
    if v := os.Getenv("MCP_HTTP_ADDR"); v != "" {
        return v
    }
    return ":8080"
}

// MCPAPIKeyFromEnv returns the MCP API key from MCP_API_KEY or a local default for demos.
func MCPAPIKeyFromEnv() string {
    if v := os.Getenv("MCP_API_KEY"); v != "" {
        return v
    }
    return "local-dev-mcp-key"
}
