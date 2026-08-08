package mcp

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/maxkondr/mcp-vm/server/metrics"
    "github.com/modelcontextprotocol/go-sdk/mcp"
    "github.com/sirupsen/logrus"
)

var (
    errMatchesRequired = errors.New("matches is required")
    errEndBeforeStart  = errors.New("end must be after or equal to start")
)

type BuildInfoInput struct{}

type BuildInfoOutput struct {
    Version   string `json:"version" jsonschema:"Vm server version"`
    Revision  string `json:"revision" jsonschema:"VM server revision"`
    Branch    string `json:"branch" jsonschema:"Vm server Git build branch"`
    BuildUser string `json:"buildUser" jsonschema:"build user"`
    BuildDate string `json:"buildDate" jsonschema:"VM server build date"`
    GoVersion string `json:"goVersion" jsonschema:"VM server Go build version"`
}

// BuildInfo is an MCP tool handler that returns VictoriaMetrics build metadata.
func BuildInfo(ctx context.Context, req *mcp.CallToolRequest, _ BuildInfoInput) (
    *mcp.CallToolResult,
    BuildInfoOutput,
    error,
) {
    _ = req
    start := metrics.ObserveToolStart("buildInfo")
    logrus.Info("tool buildInfo request received")

    res, err := vmApi.Buildinfo(ctx)
    if err != nil {
        metrics.ObserveToolEnd("buildInfo", start, err)
        logrus.WithError(err).Error("tool buildInfo request failed")
        return nil, BuildInfoOutput{}, err
    }
    metrics.ObserveToolEnd("buildInfo", start, nil)

    logrus.WithFields(logrus.Fields{
        "version":  res.Version,
        "revision": res.Revision,
    }).Info("tool buildInfo request completed")

    return nil, BuildInfoOutput{
        Version:   res.Version,
        Revision:  res.Revision,
        Branch:    res.Branch,
        BuildUser: res.BuildUser,
        BuildDate: res.BuildDate,
        GoVersion: res.GoVersion,
    }, nil
}

type SeriesInput struct {
    Matches []string `json:"match" jsonschema:"Series selectors, e.g. up or up{job=\"node\"}"`
    Start   string   `json:"start,omitempty" jsonschema:"Start time (RFC3339); default: now-1h"`
    End     string   `json:"end,omitempty" jsonschema:"End time (RFC3339); default: now"`
}

type SeriesOutput struct {
    Series   []map[string]string `json:"series" jsonschema:"Series label sets"`
    Warnings []string            `json:"warnings,omitempty" jsonschema:"Prometheus API warnings"`
}

// Series is an MCP tool handler that queries time series label sets from VictoriaMetrics.
func Series(ctx context.Context, req *mcp.CallToolRequest, in SeriesInput) (
    *mcp.CallToolResult,
    SeriesOutput,
    error,
) {
    _ = req
    startTime := metrics.ObserveToolStart("series")
    logrus.WithFields(logrus.Fields{
        "matches": len(in.Matches),
        "start":   in.Start,
        "end":     in.End,
    }).Info("tool series request received")

    if len(in.Matches) == 0 {
        metrics.ObserveToolEnd("series", startTime, errMatchesRequired)
        logrus.Error("tool series validation failed: matches is required")
        return nil, SeriesOutput{}, errMatchesRequired
    }

    end := time.Now()
    start := end.Add(-1 * time.Hour)

    var err error
    if in.End != "" {
        end, err = time.Parse(time.RFC3339, in.End)
        if err != nil {
            metrics.ObserveToolEnd("series", startTime, err)
            logrus.WithError(err).WithField("end", in.End).Error("tool series validation failed: invalid end time")
            return nil, SeriesOutput{}, fmt.Errorf("invalid end time: %w", err)
        }
    }
    if in.Start != "" {
        start, err = time.Parse(time.RFC3339, in.Start)
        if err != nil {
            metrics.ObserveToolEnd("series", startTime, err)
            logrus.WithError(err).WithField("start", in.Start).Error("tool series validation failed: invalid start time")
            return nil, SeriesOutput{}, fmt.Errorf("invalid start time: %w", err)
        }
    }
    if end.Before(start) {
        metrics.ObserveToolEnd("series", startTime, errEndBeforeStart)
        logrus.WithFields(logrus.Fields{"start": start, "end": end}).Error("tool series validation failed: end before start")
        return nil, SeriesOutput{}, errEndBeforeStart
    }

    res, warnings, err := vmApi.Series(ctx, in.Matches, start, end)
    if err != nil {
        metrics.ObserveToolEnd("series", startTime, err)
        logrus.WithError(err).Error("tool series request failed")
        return nil, SeriesOutput{}, err
    }

    series := make([]map[string]string, 0, len(res))
    for _, labelSet := range res {
        item := make(map[string]string, len(labelSet))
        for k, v := range labelSet {
            item[string(k)] = string(v)
        }
        series = append(series, item)
    }

    outWarnings := make([]string, 0, len(warnings))
    for _, w := range warnings {
        outWarnings = append(outWarnings, string(w))
    }

    logrus.WithFields(logrus.Fields{
        "series_count":   len(series),
        "warnings_count": len(outWarnings),
    }).Info("tool series request completed")
    metrics.ObserveToolEnd("series", startTime, nil)

    return nil, SeriesOutput{
        Series:   series,
        Warnings: outWarnings,
    }, nil
}

type MetadataInput struct {
    Metric string `json:"metric,omitempty" jsonschema:"Metric name to filter metadata (optional)"`
    Limit  string `json:"limit,omitempty" jsonschema:"Maximum number of metadata entries per metric (optional)"`
}

type MetadataEntryOutput struct {
    Type string `json:"type" jsonschema:"Metric type"`
    Help string `json:"help" jsonschema:"Metric help text"`
    Unit string `json:"unit" jsonschema:"Metric unit"`
}

type MetadataOutput struct {
    Metadata map[string][]MetadataEntryOutput `json:"metadata" jsonschema:"Metadata grouped by metric name"`
}

// Metadata is an MCP tool handler that returns metric metadata from VictoriaMetrics.
func Metadata(ctx context.Context, req *mcp.CallToolRequest, in MetadataInput) (
    *mcp.CallToolResult,
    MetadataOutput,
    error,
) {
    _ = req
    start := metrics.ObserveToolStart("metadata")
    logrus.WithFields(logrus.Fields{
        "metric": in.Metric,
        "limit":  in.Limit,
    }).Info("tool metadata request received")

    res, err := vmApi.Metadata(ctx, in.Metric, in.Limit)
    if err != nil {
        metrics.ObserveToolEnd("metadata", start, err)
        logrus.WithError(err).Error("tool metadata request failed")
        return nil, MetadataOutput{}, err
    }

    out := make(map[string][]MetadataEntryOutput, len(res))
    for metricName, entries := range res {
        item := make([]MetadataEntryOutput, 0, len(entries))
        for _, meta := range entries {
            item = append(item, MetadataEntryOutput{
                Type: string(meta.Type),
                Help: meta.Help,
                Unit: meta.Unit,
            })
        }
        out[metricName] = item
    }

    logrus.WithField("metrics_count", len(out)).Info("tool metadata request completed")
    metrics.ObserveToolEnd("metadata", start, nil)

    return nil, MetadataOutput{Metadata: out}, nil
}
