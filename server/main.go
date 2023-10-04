package server

import (
	//    "bytes"
	"context"
	//    "encoding/json"
	"errors"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	collectorv1 "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	// logv1 "go.opentelemetry.io/proto/otlp/logs/v1"
	// resv1 "go.opentelemetry.io/proto/otlp/resource/v1"
	// commonv1 "go.opentelemetry.io/proto/otlp/common/v1"
)

type server struct {
	collectorv1.UnimplementedLogsServiceServer
}

var err error

func (s *server) Export(ctx context.Context, in *collectorv1.ExportLogsServiceRequest) (*collectorv1.ExportLogsServiceResponse, error) {
	var alerts []SnoozeAlertV1

	rls := in.GetResourceLogs()
	if rls == nil {
    log.Error("Failed to serve log because it has no resource:", in)
    logCount.WithLabelValues("error").Inc()
		return nil, errors.New("Failed to serve log")
	}
	// For each resource log
	for _, rl := range rls {
		resource := rl.GetResource()
		sls := rl.GetScopeLogs()
		// For each scope log
		for _, sl := range sls {
			scope := sl.GetScope()
			lrs := sl.GetLogRecords()
			// For each record log
			for _, lr := range lrs {
				alert := convertAlert(resource, scope, lr)
				alerts = append(alerts, alert)
			}
		}
	}

  alertCount := float64(len(alerts))

  err = Snooze.send(alerts)
  if err != nil {
    logCount.WithLabelValues("error").Add(alertCount)
    return nil, err
  }

  logCount.WithLabelValues("ok").Add(alertCount)
	return &collectorv1.ExportLogsServiceResponse{}, nil
}

func initLogging() {
  var cll = Config.LogLevel
  if cll == "" {
    cll = "debug"
  }
  ll, err := log.ParseLevel(cll)
  if err != nil {
    log.Fatal("Unsupported log level:", Config.LogLevel)
  }
  log.SetLevel(ll)
  log.Debug("Log level set to:", ll)
}

func Run() {
  log.Infof("Starting snooze-otlp %s-%s", Version, Commit)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", Config.GrpcListeningPort))
	if err != nil {
		log.Fatalf("Failed to listen to port %d: %v", 4317, err)
	}

  if Config.PrometheusEnable {
    go serveMetrics()
  }

	s := grpc.NewServer()
	collectorv1.RegisterLogsServiceServer(s, &server{})
  healthv1.RegisterHealthServer(s, health.NewServer())
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

func InitServer() {
  initConfig()
  initLogging()
  initSnooze()
}

func init() {
  InitServer()
}
