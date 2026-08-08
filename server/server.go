package server

import (
    mcpVm "github.com/maxkondr/mcp-vm/server/mcp"
    "github.com/modelcontextprotocol/go-sdk/mcp"
)

// NewMCPServer builds and configures the MCP server with all tools and prompts.
func NewMCPServer() *mcp.Server {
    server := mcp.NewServer(&mcp.Implementation{Name: "VictoriaMetrics MCP", Version: "v1.0.0"}, nil)

    mcpVm.RegisterTools(server)
    mcpVm.RegisterPrompts(server)

    return server
}
