package mcp

import "github.com/modelcontextprotocol/go-sdk/mcp"

// RegisterTools registers all VictoriaMetrics-backed MCP tools on the server.
func RegisterTools(server *mcp.Server) {
    mcp.AddTool(server, &mcp.Tool{Name: "buildInfo", Description: "Get VictoriaMetrics build info"}, BuildInfo)
    mcp.AddTool(server, &mcp.Tool{Name: "series", Description: "Find series by label matchers"}, Series)
    mcp.AddTool(server, &mcp.Tool{Name: "metadata", Description: "Get metric metadata"}, Metadata)
}

// RegisterPrompts registers all prompt templates exposed by this MCP server.
func RegisterPrompts(server *mcp.Server) {
    server.AddPrompt(&mcp.Prompt{
        Name:        "vm_buildinfo",
        Description: "Prompt template for buildInfo tool",
    }, BuildInfoPrompt)

    server.AddPrompt(&mcp.Prompt{
        Name:        "vm_series",
        Description: "Prompt template for series tool",
        Arguments: []*mcp.PromptArgument{
            {
                Name:        "match",
                Description: "PromQL series selector (for example, up or up{job=\"node\"})",
                Required:    true,
            },
            {
                Name:        "start",
                Description: "RFC3339 start timestamp",
            },
            {
                Name:        "end",
                Description: "RFC3339 end timestamp",
            },
        },
    }, SeriesPrompt)

    server.AddPrompt(&mcp.Prompt{
        Name:        "vm_metadata",
        Description: "Prompt template for metadata tool",
        Arguments: []*mcp.PromptArgument{
            {
                Name:        "metric",
                Description: "Optional metric name to filter metadata",
            },
            {
                Name:        "limit",
                Description: "Optional max metadata entries per metric",
            },
        },
    }, MetadataPrompt)
}
