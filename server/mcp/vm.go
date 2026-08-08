package mcp

import (
    "net/http"

    "github.com/prometheus/client_golang/api"
    v1 "github.com/prometheus/client_golang/api/prometheus/v1"
    "github.com/sirupsen/logrus"
)

var (
    vmApi v1.API
)

// InitVictoriaMetricsApi initializes the shared VictoriaMetrics/Prometheus API client.
func InitVictoriaMetricsApi(vmAddress string) {
    client, err := api.NewClient(api.Config{
        Client:  http.DefaultClient,
        Address: vmAddress,
    })
    if err != nil {
        logrus.WithError(err).Fatal("failed to create VictoriaMetrics client")
    }
    vmApi = v1.NewAPI(client)
}
