package metrics

import (
    "net/http"
    "strconv"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    mcpToolCallsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "mcp_tool_calls_total",
        Help: "Total number of MCP tool calls.",
    }, []string{"tool"})
    mcpToolErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "mcp_tool_errors_total",
        Help: "Total number of MCP tool call errors.",
    }, []string{"tool"})
    mcpToolDurationSeconds = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "mcp_tool_duration_seconds",
        Help:    "Latency of MCP tool calls in seconds.",
        Buckets: prometheus.DefBuckets,
    }, []string{"tool"})

    mcpHTTPRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "mcp_http_requests_total",
        Help: "Total number of incoming HTTP requests.",
    }, []string{"method", "path", "status"})
    mcpHTTPRequestsErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "mcp_http_request_errors_total",
        Help: "Total number of HTTP requests finished with 4xx/5xx status codes.",
    }, []string{"method", "path", "status"})
    mcpHTTPRequestsDurationSeconds = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "mcp_http_request_duration_seconds",
        Help:    "Latency of incoming HTTP requests in seconds.",
        Buckets: prometheus.DefBuckets,
    }, []string{"method", "path", "status"})
)

type statusRecorder struct {
    http.ResponseWriter
    statusCode int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
    r.statusCode = statusCode
    r.ResponseWriter.WriteHeader(statusCode)
}

// Write preserves status code tracking for implicit 200 writes.
func (r *statusRecorder) Write(p []byte) (int, error) {
    return r.ResponseWriter.Write(p)
}

// Flush forwards streaming flushes required by streamable MCP responses.
func (r *statusRecorder) Flush() {
    if f, ok := r.ResponseWriter.(http.Flusher); ok {
        f.Flush()
    }
}

// Unwrap allows http.ResponseController to reach the original writer capabilities.
func (r *statusRecorder) Unwrap() http.ResponseWriter {
    return r.ResponseWriter
}

// ObserveToolStart increments per-tool call count and returns start time for latency tracking.
func ObserveToolStart(tool string) time.Time {
    mcpToolCallsTotal.WithLabelValues(tool).Inc()
    return time.Now()
}

// ObserveToolEnd records per-tool latency and increments error count when err is non-nil.
func ObserveToolEnd(tool string, start time.Time, err error) {
    if err != nil {
        mcpToolErrorsTotal.WithLabelValues(tool).Inc()
    }
    mcpToolDurationSeconds.WithLabelValues(tool).Observe(time.Since(start).Seconds())
}

// WithHTTPMetrics records request count, latency, and error status for incoming HTTP traffic.
func WithHTTPMetrics(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        recorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}

        next.ServeHTTP(recorder, r)

        status := strconv.Itoa(recorder.statusCode)
        mcpHTTPRequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
        mcpHTTPRequestsDurationSeconds.WithLabelValues(r.Method, r.URL.Path, status).Observe(time.Since(start).Seconds())
        if recorder.statusCode >= http.StatusBadRequest {
            mcpHTTPRequestsErrorsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
        }
    })
}
