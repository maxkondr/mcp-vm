package config

import "os"

const ApiKeyHeader = "X-API-Key"

// VmAddressFromEnv returns the VictoriaMetrics base URL from VM_ADDRESS or a local default.
func VmAddressFromEnv() string {
    if v := os.Getenv("VM_ADDRESS"); v != "" {
        return v
    }
    return "http://localhost:8428"
}

// HttpAddressFromEnv returns the HTTP listen address from MCP_HTTP_ADDR or a local default.
func HttpAddressFromEnv() string {
    if v := os.Getenv("MCP_HTTP_ADDR"); v != "" {
        return v
    }
    return ":8080"
}

// McpAPIKeyFromEnv returns the MCP API key from MCP_API_KEY or a local default for demos.
func McpAPIKeyFromEnv() string {
    if v := os.Getenv("MCP_API_KEY"); v != "" {
        return v
    }
    return "local-dev-mcp-key"
}
