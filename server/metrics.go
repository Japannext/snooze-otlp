package server

import (
  "fmt"
  "net/http"

  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promauto"
  "github.com/prometheus/client_golang/prometheus/promhttp"
)

var logCount = promauto.NewCounterVec(prometheus.CounterOpts{
    Name: "snooze_otlp_log_count",
    Help: "A counter for the number of log processed. Aggregated by status (ok|error)",
  },
    []string{"status"},
  )

func serveMetrics() {

  http.Handle("/metrics", promhttp.Handler())
  port := fmt.Sprintf(":%d", Config.PrometheusPort)
  http.ListenAndServe(port, nil)

}
