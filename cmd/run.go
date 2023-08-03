package cmd

import (
	"github.com/spf13/cobra"

	"github.com/japannext/snooze-otlp/server"
)

var runCmd = &cobra.Command{
  Use: "run",
  Short: "Run the snooze-otlp service",
  Long: `Run the snooze-otlp service. This will listen on Opentelemetry port
for opentelemetry gRPC format. Upon receiving a log at /v1/logs, it will
send it to snooze in the appropriate format (snooze v1)`,
  Run: func(cmd *cobra.Command, args []string) {
    server.Run()
  },
}
