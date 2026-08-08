package mcp

import (
    "context"
    "fmt"

    "github.com/modelcontextprotocol/go-sdk/mcp"
)

// BuildInfoPrompt returns a prompt template that instructs clients to call the buildInfo tool.
func BuildInfoPrompt(_ context.Context, _ *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
    return &mcp.GetPromptResult{
        Description: "Prompt for VictoriaMetrics build information",
        Messages: []*mcp.PromptMessage{
            {
                Role:    mcp.Role("user"),
                Content: &mcp.TextContent{Text: "Call the buildInfo tool and provide the result as a short, readable summary."},
            },
        },
    }, nil
}

// SeriesPrompt returns a prompt template that instructs clients to call the series tool.
func SeriesPrompt(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
    match := "up"
    start := ""
    end := ""

    if req != nil && req.Params != nil && req.Params.Arguments != nil {
        if req.Params.Arguments["match"] != "" {
            match = req.Params.Arguments["match"]
        }
        start = req.Params.Arguments["start"]
        end = req.Params.Arguments["end"]
    }

    text := fmt.Sprintf(
        "Call the series tool with match=[%q]. Include start/end only if provided, then summarize the labels returned.",
        match,
    )
    if start != "" || end != "" {
        text = fmt.Sprintf(
            "Call the series tool with match=[%q], start=%q, end=%q, then summarize the labels returned.",
            match,
            start,
            end,
        )
    }

    return &mcp.GetPromptResult{
        Description: "Prompt for VictoriaMetrics series lookup",
        Messages: []*mcp.PromptMessage{
            {
                Role:    mcp.Role("user"),
                Content: &mcp.TextContent{Text: text},
            },
        },
    }, nil
}

// MetadataPrompt returns a prompt template that instructs clients to call the metadata tool.
func MetadataPrompt(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
    metric := ""
    limit := ""

    if req != nil && req.Params != nil && req.Params.Arguments != nil {
        metric = req.Params.Arguments["metric"]
        limit = req.Params.Arguments["limit"]
    }

    text := "Call the metadata tool and summarize metric metadata grouped by metric name."
    if metric != "" || limit != "" {
        text = fmt.Sprintf(
            "Call the metadata tool with metric=%q and limit=%q, then summarize results grouped by metric name.",
            metric,
            limit,
        )
    }

    return &mcp.GetPromptResult{
        Description: "Prompt for VictoriaMetrics metadata lookup",
        Messages: []*mcp.PromptMessage{
            {
                Role:    mcp.Role("user"),
                Content: &mcp.TextContent{Text: text},
            },
        },
    }, nil
}
