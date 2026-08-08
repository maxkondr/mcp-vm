package auth

import (
    "crypto/subtle"
    "net/http"

    "github.com/maxkondr/mcp-vm/server/config"
    "github.com/sirupsen/logrus"
)

// WithAPIKeyAuth enforces API key authentication on MCP HTTP requests.
func WithAPIKeyAuth(next http.Handler, apiKey string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if subtle.ConstantTimeCompare([]byte(r.Header.Get(config.APIKeyHeader)), []byte(apiKey)) != 1 {
            logrus.WithFields(logrus.Fields{
                "remote_addr": r.RemoteAddr,
                "path":        r.URL.Path,
            }).Warn("MCP request rejected: invalid API key")
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}
