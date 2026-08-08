package main

import (
    "crypto/tls"
    "errors"
    "net/http"
    "time"

    "github.com/maxkondr/mcp-vm/server"
    "github.com/maxkondr/mcp-vm/server/auth"
    "github.com/maxkondr/mcp-vm/server/config"
    "github.com/maxkondr/mcp-vm/server/metrics"
    "github.com/modelcontextprotocol/go-sdk/mcp"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/sirupsen/logrus"
)

func main() {
    logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
    logrus.SetLevel(logrus.InfoLevel)

    tr := http.DefaultTransport.(*http.Transport).Clone()
    tr.TLSClientConfig = &tls.Config{
        InsecureSkipVerify: true, //nolint:gosec // only for local/dev testing
    }

    // Apply globally to default client.
    http.DefaultClient.Transport = tr
    http.DefaultClient.Timeout = 10 * time.Second

    vmAddress := config.VmAddressFromEnv()
    logrus.WithField("vm_address", vmAddress).Info("initializing VictoriaMetrics client")

    // MCP
    server := server.NewMCPServer()
    handler := mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server {
        return server
    }, &mcp.StreamableHTTPOptions{
        SessionTimeout: 5 * time.Minute,
    })
    authenticatedHandler := auth.WithAPIKeyAuth(handler, config.McpAPIKeyFromEnv())

    mux := http.NewServeMux()
    mux.Handle("/mcp", authenticatedHandler)
    mux.Handle("/metrics", promhttp.Handler())
    instrumentedMux := metrics.WithHTTPMetrics(mux)

    httpServer := &http.Server{
        Addr:              config.HttpAddressFromEnv(),
        Handler:           instrumentedMux,
        ReadHeaderTimeout: 5 * time.Second,
    }

    logrus.WithFields(logrus.Fields{
        "addr":             httpServer.Addr,
        "path":             "/mcp",
        "metrics_endpoint": "/metrics",
        "auth_header":      config.ApiKeyHeader,
        "api_key_enabled":  true,
    }).Info("starting MCP streamable HTTP server")

    err := httpServer.ListenAndServe()
    if err != nil && !errors.Is(err, http.ErrServerClosed) {
        logrus.WithError(err).Fatal("MCP HTTP server failed")
    }
}
