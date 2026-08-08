package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// mcpEndpointFromEnv returns the MCP Streamable HTTP endpoint from env or local default.
func mcpEndpointFromEnv() string {
	if v := os.Getenv("MCP_SERVER_URL"); v != "" {
		return v
	}
	return "http://localhost:8080/mcp"
}

// mcpAPIKeyFromEnv returns the MCP API key from MCP_API_KEY or a local default for demos.
func mcpAPIKeyFromEnv() string {
	if v := os.Getenv("MCP_API_KEY"); v != "" {
		return v
	}
	return "local-dev-mcp-key"
}

// apiKeyRoundTripper injects an API key header in outgoing HTTP requests.
type apiKeyRoundTripper struct {
	base   http.RoundTripper
	apiKey string
}

// RoundTrip sends the request with MCP API key authentication header.
func (t *apiKeyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	cloned := req.Clone(req.Context())
	cloned.Header = req.Header.Clone()
	cloned.Header.Set("X-API-Key", t.apiKey)

	return t.base.RoundTrip(cloned)
}

// printToolResult logs content and structured output returned by a tool call.
func printToolResult(res *mcp.CallToolResult) {
	for _, c := range res.Content {
		switch v := c.(type) {
		case *mcp.TextContent:
			log.Print(v.Text)
		default:
			b, err := json.Marshal(v)
			if err != nil {
				log.Printf("failed to marshal tool content: %v", err)
				continue
			}
			log.Print(string(b))
		}
	}

	if len(res.Content) == 0 && res.StructuredContent != nil {
		b, err := json.MarshalIndent(res.StructuredContent, "", "  ")
		if err != nil {
			log.Printf("failed to marshal structured content: %v", err)
			return
		}
		log.Printf("structuredContent: %s", string(b))
	}
}

func main() {
	ctx := context.Background()

	// Create a new client, with no features.
	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)
	apiKey := mcpAPIKeyFromEnv()

	// Connect to the server over streamable HTTP transport.
	transport := &mcp.StreamableClientTransport{
		Endpoint: mcpEndpointFromEnv(),
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &apiKeyRoundTripper{
				base:   http.DefaultTransport,
				apiKey: apiKey,
			},
		},
	}
	log.Printf("Connecting to MCP server at %s with API key auth", transport.Endpoint)

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// Call buildInfo tool on the server.
	params := &mcp.CallToolParams{
		Name:      "buildInfo",
		Arguments: map[string]any{},
	}
	log.Print("Calling buildInfo tool...")
	res, err := session.CallTool(ctx, params)
	if err != nil {
		log.Fatalf("CallTool failed: %v", err)
	}
	if res.IsError {
		log.Fatalf("tool failed: %v", res.GetError())
	}
	printToolResult(res)

	// Call series tool on the server.
	params = &mcp.CallToolParams{
		Name:      "series",
		Arguments: map[string]any{"match": []string{"go_sql_connections_max_open"}},
	}
	log.Print("Calling series tool...")
	res, err = session.CallTool(ctx, params)
	if err != nil {
		log.Fatalf("CallTool failed: %v", err)
	}
	if res.IsError {
		log.Fatalf("tool failed: %v", res.GetError())
	}
	printToolResult(res)

	// Call metadata tool on the server.
	params = &mcp.CallToolParams{
		Name:      "metadata",
		Arguments: map[string]any{"metric": "go_sql_connections_max_open"},
	}
	log.Print("Calling metadata tool...")
	res, err = session.CallTool(ctx, params)
	if err != nil {
		log.Fatalf("CallTool failed: %v", err)
	}
	if res.IsError {
		log.Fatalf("tool failed: %v", res.GetError())
	}
	printToolResult(res)
}
